package service

import (
	"bufio"
	"context"
	"doo-store/backend/config"
	"doo-store/backend/constant"
	"doo-store/backend/core/dto"
	"doo-store/backend/core/dto/request"
	"doo-store/backend/core/dto/response"
	"doo-store/backend/core/model"
	"doo-store/backend/core/repo"
	"doo-store/backend/task"
	"doo-store/backend/utils/common"
	"doo-store/backend/utils/compose"
	"doo-store/backend/utils/docker"
	e "doo-store/backend/utils/error"
	"doo-store/backend/utils/nginx"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"gorm.io/gorm"
)

type AppService struct {
}

type IAppService interface {
	AppPage(ctx dto.ServiceContext, req request.AppSearch) (*dto.PageResult, error)
	AppDetailByKey(ctx dto.ServiceContext, key string) (*response.AppDetail, error)
	AppInstall(ctx dto.ServiceContext, req request.AppInstall) error
	AppInstallOperate(ctx dto.ServiceContext, req request.AppInstalledOperate) error
	AppUnInstall(ctx dto.ServiceContext, req request.AppUnInstall) error
	AppInstalledPage(ctx dto.ServiceContext, req request.AppInstalledSearch) (*dto.PageResult, error)
	Params(ctx dto.ServiceContext, id int64) (any, error)
	UpdateParams(ctx dto.ServiceContext, req request.AppInstall) (any, error)
	AppTags(ctx dto.ServiceContext) ([]*model.Tag, error)
	GetLogs(ctx dto.ServiceContext, req request.AppLogsSearch) (any, error)
	Upload(ctx dto.ServiceContext, req request.PluginUpload) error
}

func NewIAppService() IAppService {
	return &AppService{}
}

type AppInstallProcess struct {
	ctx           dto.ServiceContext
	app           *model.App
	appDetail     *model.AppDetail
	appInstalled  *model.AppInstalled
	appKey        string
	containerName string
	envContent    string
	req           request.AppInstall
	ipAddress     string
	client        docker.Client
}

func NewAppInstallProcess(ctx dto.ServiceContext, req request.AppInstall) *AppInstallProcess {
	return &AppInstallProcess{
		ctx: ctx,
		req: req,
	}
}

func (p *AppInstallProcess) Check() error {
	var err error
	p.app, err = repo.App.Where(repo.App.Key.Eq(p.req.Key)).First()
	if err != nil {
		log.Info("Error query app")
		return errors.New("获取插件信息失败")
	}
	// 检测版本
	dootaskService := NewIDootaskService()
	versionInfoResp, err := dootaskService.GetVersoinInfo()
	if err != nil {
		return errors.New("获取版本信息失败")
	}
	check, err := versionInfoResp.CheckVersion(p.app.DependsVersion)
	if err != nil {
		log.Info("检测版本失败", err)
		return errors.New("检查依赖版本失败")
	}
	// 依赖版本不符合要求
	if !check {
		return e.WithMap(p.ctx.C, constant.ErrPluginVersionNotSupport, map[string]interface{}{
			"detail": p.app.DependsVersion,
		}, nil)
	}

	// 判断是否已安装
	p.appInstalled, err = repo.AppInstalled.
		Select(repo.AppInstalled.ID, repo.AppInstalled.AppID).
		Where(repo.AppInstalled.AppID.Eq(p.app.ID)).
		First()

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return errors.New("安装失败")
		}
	}
	if p.appInstalled != nil {
		return errors.New("无需重复安装")
	}

	p.appDetail, err = repo.AppDetail.Select(
		repo.AppDetail.ID,
		repo.AppDetail.AppID,
		repo.AppDetail.Repo,
		repo.AppDetail.Version,
		repo.AppDetail.Params,
		repo.AppDetail.DependsVersion,
		repo.AppDetail.NginxConfig,
	).Where(repo.AppDetail.AppID.Eq(p.app.ID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return errors.New("安装失败")
	}
	return nil
}

func (p *AppInstallProcess) DHCP() error {
	var err error
	p.client, err = docker.NewClient()
	if err != nil {
		return err
	}
	client := p.client.GetClient()

	containers, err := client.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})
	if err != nil {
		log.Info("获取容器列表失败")
		return err
	}
	networks, err := client.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		log.Info("获取容器网络列表失败")
		return err
	}
	networkID := ""
	for _, network := range networks {
		if network.Name == config.EnvConfig.GetNetworkName() {
			networkID = network.ID
			break
		}
	}
	if networkID == "" {
		return errors.New("网络不存在")
	}
	var usedIPs []string
	fmt.Println("容器数量", len(containers))
	for _, container := range containers {
		// 获取容器的网络设置
		for _, network := range container.NetworkSettings.Networks {
			fmt.Println("容器使用的网络:", network.NetworkID)
			if network.NetworkID == networkID {
				// dootask-networks-5d2bc8
				usedIPs = append(usedIPs, network.IPAddress)
			}
		}
	}
	allInstalled, _ := repo.AppInstalled.Select(repo.AppInstalled.IpAddress).Find()
	for _, installed := range allInstalled {
		usedIPs = append(usedIPs, installed.IpAddress)
	}

	// 获取一个未使用的 IP 地址
	getAvailableIP := func(usedIPs []string) (string, error) {

		ipParts := strings.Split(config.EnvConfig.IP_START, ".")
		if len(ipParts) != 4 {
			return "", errors.New("IP 地址格式不正确")
		}
		ip, err := strconv.Atoi(ipParts[3])
		if err != nil {
			return "", err
		}
		for i := ip; i < ip+config.EnvConfig.IP_COUNT; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%d", ipParts[0], ipParts[1], ipParts[2], i)
			if !common.InArray(ip, usedIPs) {
				return ip, nil
			}
		}

		return "", errors.New("没有可用的 IP 地址")
	}

	p.ipAddress, err = getAvailableIP(usedIPs)
	if err != nil {
		return err
	}
	return nil
}

// ValidateParam 验证参数
func (p *AppInstallProcess) ValidateParam() error {
	var err error
	// 检测 docker-compose 文件
	err = compose.Check(p.req.DockerCompose)
	if err != nil {
		log.Info("DockerCompose 内容未通过检测", err)
		return err
	}

	p.appKey = config.EnvConfig.APP_PREFIX + p.app.Key
	// 创建工作目录
	workspaceDir := path.Join(constant.AppInstallDir, p.appKey)
	err = createDir(workspaceDir)
	if err != nil {
		log.Info("Error create dir", err)
		return err
	}

	// 容器名称
	p.containerName = config.EnvConfig.GetDefaultContainerName(p.app.Key)

	paramJson, err := json.Marshal(p.req.Params)
	if err != nil {
		return err
	}

	params := response.AppParams{}
	err = common.StrToStruct(p.appDetail.Params, &params)
	if err != nil {
		log.Debug("解析参数失败", err)
		return errors.New("解析插件参数失败")
	}
	for _, param := range params.FormFields {
		if param.Required {
			if _, exists := p.req.Params[param.EnvKey]; !exists {
				return errors.New("缺少必填参数 " + param.EnvKey)
			}
		}
	}

	// 资源限制
	p.req.Params[constant.CPUS] = p.req.CPUS
	p.req.Params[constant.MemoryLimit] = p.req.MemoryLimit
	var envJson string
	p.envContent, envJson, err = docker.GenEnv(p.appKey, p.containerName, p.ipAddress, p.req.Params, false)
	if err != nil {
		return err
	}
	p.appInstalled = &model.AppInstalled{
		Name:          p.containerName,
		AppID:         p.app.ID,
		AppDetailID:   p.appDetail.ID,
		Class:         p.app.Class,
		Repo:          p.appDetail.Repo,
		Version:       p.appDetail.Version,
		Params:        string(paramJson),
		Env:           envJson,
		DockerCompose: p.req.DockerCompose,
		Key:           p.app.Key,
		Status:        constant.Installing,
		IpAddress:     p.ipAddress,
	}
	// 更新插件状态
	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		_, err = repo.Use(tx).App.Where(repo.App.ID.Eq(p.appInstalled.AppID)).Update(repo.App.Status, constant.AppInUse)
		if err != nil {
			return err
		}
		err = repo.Use(tx).AppInstalled.Create(p.appInstalled)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *AppInstallProcess) Install() error {
	var err error
	if p.appInstalled == nil {
		return errors.New("安装失败")
	}
	err = appUp(p.appInstalled, p.envContent)
	if err != nil {
		log.Info("启动失败", err)
		return err
	}
	return nil
}

// AddNginx 添加Nginx配置
// 插件安装的时候，需要向Nginx添加一个配置，如果添加配置失败，会将插件停止
func (p *AppInstallProcess) AddNginx() error {
	client, err := docker.NewClient()
	if err != nil {
		return err
	}
	port, err := client.GetImageFirstExposedPortByName(fmt.Sprintf("%s:%s", p.appDetail.Repo, p.appDetail.Version))
	if err != nil {
		return err
	}
	if p.appDetail.NginxConfig != "" || port != 0 {
		err = nginx.AddLocation(p.appDetail.NginxConfig, p.app.Key, p.containerName, port)

		if err != nil {
			log.Info("添加Nginx配置失败", err)

			std, err := compose.Operate(docker.GetComposeFile(p.appKey), "stop")
			if err != nil {
				log.Info("Error docker compose operate", std)
			}
			_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(p.appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
			return err
		}
	}
	return nil
}

func (*AppService) AppPage(ctx dto.ServiceContext, req request.AppSearch) (*dto.PageResult, error) {
	var query repo.IAppDo
	query = repo.App.Order(repo.App.Sort.Desc())
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 9
	} else if req.PageSize > 1000 {
		req.PageSize = 1000
	}
	if req.Name != "" {
		query = query.Where(repo.App.Name.Like(fmt.Sprintf("%%%s%%", req.Name)))
	}
	if req.Class != "" {
		query = query.Where(repo.App.Class.Eq(req.Class))
	}
	if req.ID != 0 {
		query = query.Where(repo.App.ID.Eq(req.ID))
	}
	if req.Description != "" {
		query = query.Where(repo.App.Description.Like(fmt.Sprintf("%%%s%%", req.Description)))
	}
	result, count, err := query.FindByPage((req.Page-1)*req.PageSize, req.PageSize)

	if err != nil {
		return nil, err
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) AppDetailByKey(ctx dto.ServiceContext, key string) (*response.AppDetail, error) {

	app, err := repo.App.Where(repo.App.Key.Eq(key)).First()
	if err != nil {
		return nil, err
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.AppID.Eq(app.ID)).First()
	if err != nil {
		return nil, err
	}
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		return nil, err
	}
	resp := &response.AppDetail{
		AppDetail: *appDetail,
		Params:    params,
	}

	return resp, nil
}

// AppInstall 插件安装
func (*AppService) AppInstall(ctx dto.ServiceContext, req request.AppInstall) error {
	appInstallProcess := NewAppInstallProcess(ctx, req)
	if err := appInstallProcess.Check(); err != nil {
		return err
	}
	if err := appInstallProcess.DHCP(); err != nil {
		return err
	}
	if err := appInstallProcess.ValidateParam(); err != nil {
		return err
	}
	// 异步处理
	manager := task.GetGlobalManager()
	manager.AddTask(func() error {
		if err := appInstallProcess.Install(); err != nil {
			return err
		}
		if err := appInstallProcess.AddNginx(); err != nil {
			return err
		}
		return nil
	})

	return nil
}

func (*AppService) AppInstallOperate(ctx dto.ServiceContext, req request.AppInstalledOperate) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)

	supportActions := []string{"start", "stop"}
	if !common.InArray(req.Action, supportActions) {
		return errors.New("不支持的action")
	}

	if req.Action == "stop" {
		err := appStop(appInstalled)
		return err
	}

	stdout, err := compose.Operate(composeFile, req.Action)
	if err != nil {
		log.Info("Error docker compose operate")
		return err
	}
	fmt.Println(stdout)
	insertLog(appInstalled.ID, fmt.Sprintf("插件操作[%s]", req.Action), stdout)
	return nil
}

func (*AppService) AppUnInstall(ctx dto.ServiceContext, req request.AppUnInstall) error {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.Key.Eq(req.Key)).First()
	if err != nil {
		return err
	}
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	err = repo.DB.Transaction(func(tx *gorm.DB) error {
		_, err = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Delete()
		if err != nil {
			log.Info("删除插件失败", err)
			return err
		}
		_, err = repo.Use(tx).App.Where(repo.App.ID.Eq(appInstalled.AppID)).Update(repo.App.Status, constant.AppUnused)
		if err != nil {
			log.Info("更新插件状态失败", err)
			return err
		}
		stdout, err := compose.Down(composeFile)
		if err != nil {
			log.Info("Error docker compose down")
			return err
		}
		fmt.Println(stdout)
		return err
	})
	if err != nil {
		log.Info("插件卸载失败", err)
		return errors.New("插件卸载失败")
	}

	nginx.RemoveLocation(appInstalled.Key)
	// 删除compose目录
	_ = os.RemoveAll(fmt.Sprintf("%s/%s", constant.AppInstallDir, appKey))

	return nil
}

func (*AppService) AppInstalledPage(ctx dto.ServiceContext, req request.AppInstalledSearch) (*dto.PageResult, error) {

	query := repo.AppInstalled.Join(repo.App, repo.App.ID.EqCol(repo.AppInstalled.AppID))
	if req.Class != "" {
		query = query.Where(repo.AppInstalled.Class.Eq(req.Class))
	}
	if req.Name != "" {
		query = query.Where(repo.App.Name.Like(fmt.Sprintf("%%%s%%", req.Name)))
	}
	if req.Description != "" {
		query = query.Where(repo.App.Description.Like(fmt.Sprintf("%%%s%%", req.Description)))
	}

	result := []map[string]any{}
	count, err := query.Select(repo.AppInstalled.ALL, repo.App.Icon, repo.App.Description, repo.App.Name).ScanByPage(&result, (req.Page-1)*req.PageSize, req.PageSize)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &dto.PageResult{
				Total: 0,
				Items: []interface{}{},
			}, nil
		}
		log.Info("查询已安装插件失败", err)
		return nil, errors.New("查询已安装插件失败")
	}

	pageResult := &dto.PageResult{
		Total: count,
		Items: result,
	}
	return pageResult, nil
}

func (*AppService) Params(ctx dto.ServiceContext, id int64) (any, error) {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(id)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.ID.Eq(appInstalled.AppDetailID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	// appDetail.Params
	// 解析原始参数
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		log.Info("错误解析Json", err)
		return nil, err
	}
	env := map[string]string{}
	err = json.Unmarshal([]byte(appInstalled.Env), &env)
	if err != nil {
		log.Info("解析环境变量失败", err)
		return nil, err
	}
	for _, formField := range params.FormFields {
		formField.Value = env[formField.EnvKey]
		formField.Key = formField.EnvKey
	}
	// 构建插件参数
	aParams := response.AppInstalledParamsResp{
		Params:        params.FormFields,
		DockerCompose: appInstalled.DockerCompose,
		CPUS:          env[constant.CPUS],
		MemoryLimit:   env[constant.MemoryLimit],
	}
	return aParams, nil
}

func (*AppService) UpdateParams(ctx dto.ServiceContext, req request.AppInstall) (any, error) {
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(req.InstalledId)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	appDetail, err := repo.AppDetail.Where(repo.AppDetail.ID.Eq(appInstalled.AppDetailID)).First()
	if err != nil {
		log.Info("Error query app detail", err)
		return nil, errors.New("获取安装插件信息失败")
	}
	// appDetail.Params
	// 解析原始参数
	params := response.AppParams{}
	err = common.StrToStruct(appDetail.Params, &params)
	if err != nil {
		log.Info("错误解析Json", err)
		return nil, err
	}
	// TODO 参数校验
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	containerName := appInstalled.Name
	ipAddress := appInstalled.IpAddress

	req.Params[constant.CPUS] = req.CPUS
	req.Params[constant.MemoryLimit] = req.MemoryLimit

	envContent, envJson, err := docker.GenEnv(appKey, containerName, ipAddress, req.Params, false)
	if err != nil {
		log.Info("错误生成环境变量文件", err)
		return nil, errors.New("修改参数失败")
	}
	appInstalled.Env = envJson
	paramJson, err := json.Marshal(req.Params)
	if err != nil {
		return nil, errors.New("解析参数失败")
	}
	appInstalled.Params = string(paramJson)
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(appInstalled)
	err = appRe(appInstalled, envContent)
	if err != nil {
		log.Info("重启失败", err)
		insertLog(appInstalled.ID, "插件重启", err.Error())
		return nil, errors.New("插件重启失败")
	}
	// 返回修改后的参数
	env := map[string]string{}
	err = json.Unmarshal([]byte(appInstalled.Env), &env)
	if err != nil {
		log.Info("解析环境变量失败", err)
		return nil, err
	}
	for _, formField := range params.FormFields {
		formField.Value = env[formField.EnvKey]
		formField.Key = formField.EnvKey
	}
	aParams := response.AppInstalledParamsResp{
		Params:        params.FormFields,
		DockerCompose: appInstalled.DockerCompose,
		CPUS:          req.CPUS,
		MemoryLimit:   req.MemoryLimit,
	}
	insertLog(appInstalled.ID, "插件参数修改", "")
	return aParams, nil
}

func (*AppService) AppTags(ctx dto.ServiceContext) ([]*model.Tag, error) {
	tags, err := repo.Tag.Find()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*model.Tag{}, nil
		}
		return nil, err
	}
	return tags, nil
}

func (*AppService) GetLogs(ctx dto.ServiceContext, req request.AppLogsSearch) (any, error) {
	log.Info("获取日志")
	client, err := docker.NewDockerClient()
	if err != nil {
		log.Info("获取Docker失败", err)
		return nil, err
	}
	appInstalled, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(req.Id)).First()
	if err != nil {
		log.Info("Error query app installed", err)
		return nil, errors.New("获取安装插件信息失败")
	}

	if appInstalled.Status != constant.Running {
		return nil, errors.New("插件未运行")
	}

	reader, err := client.ContainerLogs(context.Background(), appInstalled.Name, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		// Follow:     true,
		Since: req.Since,
		Until: req.Until,
		Tail:  fmt.Sprintf("%d", req.Tail),
	})
	if err != nil {
		log.Info("Error query container log", err)
		return nil, errors.New("获取日志失败")
	}
	sb := strings.Builder{}

	scanner := bufio.NewScanner(reader)

	cleanText := func(input string) string {
		// 去除所有不可见的控制字符 (ASCII 0-31 和 127)
		reControlChars := regexp.MustCompile(`[\x00-\x1F\x7F]`)
		str := reControlChars.ReplaceAllString(input, "")
		if len(str) > 2 {
			return str[1:]
		}
		return str
	}
	date := time.Now()
	go func() {

		for scanner.Scan() {
			fmt.Println(scanner.Text())
			_, err = sb.WriteString(cleanText(scanner.Text()) + "\n")
			if err != nil {
				log.Info("写入字符串时失败", err)
			}
			date = time.Now()
		}

		if err := scanner.Err(); err != nil {
			log.Info("Error scan container log", err)
			return
		}

	}()
	for {
		time.Sleep(time.Microsecond * 500)
		if time.Since(date) > time.Second {
			log.Info("距离上次读取时间超过1秒，读取结束")
			break
		}
	}
	return sb.String(), nil
}

// Upload 插件上传
func (AppService) Upload(ctx dto.ServiceContext, req request.PluginUpload) error {
	key := req.Plugin.Key
	count, err := repo.App.Where(repo.App.Key.Eq(key)).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("key已经存在")
	}
	err = repo.DB.Transaction(func(tx *gorm.DB) error {

		app := &model.App{
			Name:           req.Plugin.Name,
			Key:            req.Plugin.Key,
			Icon:           req.Plugin.Icon,
			Class:          req.Plugin.Class,
			Description:    req.Plugin.Description,
			DependsVersion: req.Plugin.DependsVersion,
			Status:         constant.AppUnused,
		}
		err := repo.Use(tx).App.Create(app)
		if err != nil {
			log.Debug(err.Error())
			return err
		}
		dockerCompose := req.DockerCompose
		// 生成默认docker-compose.yml
		if dockerCompose == "" {
			dockerCompose = req.Plugin.GenComposeFile()
		}

		err = compose.Check(dockerCompose)
		if err != nil {
			return err
		}

		nginxConfig := req.NginxConfig
		if nginxConfig == "" {
			nginxConfig = req.Plugin.GenNginxConfig()
		}

		appDetail := &model.AppDetail{
			AppID:          app.ID,
			Repo:           req.Plugin.Repo,
			Version:        req.Plugin.Version,
			DependsVersion: req.Plugin.DependsVersion,
			Params:         req.Plugin.GenParams(),
			DockerCompose:  dockerCompose,
			NginxConfig:    nginxConfig,
			Status:         constant.AppNormal,
		}
		err = repo.Use(tx).AppDetail.Create(appDetail)
		if err != nil {
			log.Debug(err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func appRe(appInstalled *model.AppInstalled, envContent string) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	_, err := compose.Down(composeFile)
	if err != nil {
		log.Info("Error docker compose down", err)
		return err
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Installing)
	// 写入docker-compose.yaml和环境文件
	composeFile, err = docker.WriteComposeFile(appKey, appInstalled.DockerCompose)
	if err != nil {
		log.Error("DockerCompose文件写入失败", err)
		return err
	}
	_, err = docker.WriteEnvFile(appKey, envContent)
	if err != nil {
		log.Error("环境变量文件写入失败", err)
		return err
	}
	stdout, err := compose.Up(composeFile)
	if err != nil {
		log.Info("Error docker compose up", stdout)
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.UpErr)
		return err
	}
	_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Running)

	return nil
}

// appUp
// envContent key=value
func appUp(appInstalled *model.AppInstalled, envContent string) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		composeFile, err := docker.WriteComposeFile(appKey, appInstalled.DockerCompose)
		log.Info("Docker容器UP,", composeFile)
		if err != nil {
			log.Info("Error WriteFile", err)
			return err
		}
		_, err = docker.WriteEnvFile(appKey, envContent)
		if err != nil {
			log.Info("Error WriteFile", err)
			return err
		}
		stdout, err := compose.Up(composeFile)
		if err != nil {
			stdout, err = docker.ParseError(stdout, err)
			log.Info("Error docker compose up:", stdout, err)
			return err
		}
		fmt.Println(stdout)
		_, err = repo.Use(tx).AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			model.AppInstalled{
				Status:  constant.Running,
				Message: "",
			},
		)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		_, _ = repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Updates(
			model.AppInstalled{
				Status:  constant.UpErr,
				Message: err.Error(),
			},
		)
		insertLog(appInstalled.ID, "插件启动", err.Error())
	} else {
		insertLog(appInstalled.ID, "插件启动", "")
	}
	return err
}

// appStop 插件停止
func appStop(appInstalled *model.AppInstalled) error {
	appKey := config.EnvConfig.APP_PREFIX + appInstalled.Key
	composeFile := docker.GetComposeFile(appKey)
	_, err := repo.AppInstalled.Where(repo.AppInstalled.ID.Eq(appInstalled.ID)).Update(repo.AppInstalled.Status, constant.Stopped)
	if err != nil {
		return err
	}
	stdout, err := compose.Stop(composeFile)
	if err != nil {
		return fmt.Errorf("error docker compose stop: %s", err.Error())
	}
	insertLog(appInstalled.ID, "插件停止", stdout)
	return nil
}

func createDir(dirPath string) error {
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			log.WithField("file", dirPath).Debug("file exists")
			return nil
		}
		return err
	}
	return nil
}

func insertLog(appInstalledId int64, prefix, content string) {
	if prefix == "" && content == "" {
		log.Info("log content is empty")
		return
	}
	err := repo.AppLog.Create(&model.AppLog{
		AppInstalledId: appInstalledId,
		Content:        fmt.Sprintf("%s-%s", prefix, content),
	})
	if err != nil {
		log.Info("Error create app log")
	}
}

// ConvertToUTF8 尝试将非 UTF-8 内容转换为 UTF-8
func ConvertToUTF8(input []byte) (string, error) {
	// 尝试使用 GBK 解码（示例，可以替换为其他编码）
	reader := transform.NewReader(strings.NewReader(string(input)), simplifiedchinese.GBK.NewDecoder())
	converted, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(converted), nil
}
