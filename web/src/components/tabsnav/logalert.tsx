/* eslint-disable react-hooks/exhaustive-deps */
/* eslint-disable @typescript-eslint/no-unused-vars */
import Codemirror, {ReactCodeMirrorRef} from "@uiw/react-codemirror";
import { javascript } from '@codemirror/lang-javascript';
import { Select, SelectItem, SelectTrigger, SelectValue, SelectContent } from "@/components/ui/select";
import { useState, useEffect, useRef } from "react";
import { Item } from "@/type.d/common";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Label } from '@/components/ui/label';
import * as http from "@/api/modules/fouceinter";
import { useTranslation } from "react-i18next";

interface AlertLogHaveProps {
    isLogOpen: boolean;
    isOpen: boolean;
    onClose: () => void;
    app: Item;
}

export function AlertLogHave({ isOpen, onClose, app }: AlertLogHaveProps) {

    const { t } = useTranslation();

    const [logInfo, setLogInfo] = useState('');
    const codemirrorRef = useRef<ReactCodeMirrorRef>(null);  //创建对 Codemirror 编辑器的引用，便于侧边滚动到最后
    const [logSearch, setLogSearch] = useState({
        modeIndex: 0, // 当前选中的时间范围索引
        tail: 10000,
    });


    // 时间范围选择项
    const timeOptions = [
        { label: 'All Log', value: 'all' }, // 获取所有日志
        { label: 'Last day', value: (Math.round((Date.now() - 24 * 60 * 60 * 1000) / 1000)) }, // 昨天
        { label: 'Last 4 hours', value: (Math.round((Date.now() - 4 * 60 * 60 * 1000) / 1000)) }, // 最近4小时
        { label: 'Last 1 hour', value: (Math.round((Date.now() - 1 * 60 * 60 * 1000) / 1000)) }, // 最近1小时
        { label: 'Last 10 minutes', value: (Math.round((Date.now() - 10 * 60 * 1000) / 1000)) }, // 最近10分钟
    ];

    useEffect(() => {
        if (!isOpen) return;

        fetchLogs(); // 当打开日志抽屉时，拉取日志
        return () => {
            // 关闭日志时的清理工作（没有 WebSocket连接，所以不需要关闭）
        };
    }, [isOpen, logSearch]); // isOpen：true 或 logSearch 更新时，会调用 fetchLogs 函数请求拉取日志

    // 使用 axios来请求日志
    const fetchLogs = async () => {
        try {
            // 获取当前选中的时间戳
            const selectedOption = timeOptions[logSearch.modeIndex];
            const logparams = {
                tail: logSearch.tail,
                since: selectedOption.value !== 'all' ? selectedOption.value : undefined, // 只传递选中的时间戳
            };
            const response = await http.getLogs(app.id, logparams)
            setLogInfo(response.data ||''); // 设置返回的日志数据
            console.log("日志 data:", response.data);

            //  设置Codemirror 的引用，滚动到最后
            if (codemirrorRef.current) {
                if(codemirrorRef.current.view){
                    const { doc } = codemirrorRef.current.view.state;
                    codemirrorRef.current.view.dispatch({
                        selection: { anchor: doc.length, head: doc.length },
                        scrollIntoView: true,
                });
            }
            }
        } catch (error) {
            console.log("Failed to fetch logs:", error);
            alert('container.fetchLogError'); // 错误提示

        }
    };

    // 更新选中的时间段的值，使用下标更新
    const handleModeChange = (index: number) => {
        setLogSearch(prev => ({ ...prev, modeIndex: index }));
    };

    return (
        <Sheet open={isOpen} onOpenChange={onClose}>
            <SheetContent className="overflow-y-auto overflow-x-hidden">
                <SheetHeader>
                    <SheetTitle className=' top-3 my-1 text-gray-700 text-xl'>Log</SheetTitle>
                </SheetHeader>
                <hr />
                <SheetDescription className='pt-3'>
                </SheetDescription>
                <div className='items-center mt-3'>
                    <div className="flex justify-between w-full">
                        <Label className='ml-1 w-1/6' >{t('时间范围')}</Label>
                        <Select value={logSearch.modeIndex.toString()} onValueChange={(value) => handleModeChange(Number(value))}>
                            <SelectTrigger className="lg:w-full md:w-52 w-28 bg-gray-200/60">
                                <SelectValue>
                                    {/* 通过下标渲染选中的时间范围的label */}
                                    {timeOptions[logSearch.modeIndex].label}
                                </SelectValue>
                            </SelectTrigger>
                            <SelectContent>
                                {timeOptions.map((item, index) => (
                                    <SelectItem key={index} value={index.toString()}>
                                        {item.label}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>

                    <div className="flex justify-between w-full mt-6">
                        {/* 条数选择 */}
                        <Label className='ml-1 w-1/6'>{t('条数')}</Label>
                        <Select 
                            value={logSearch.tail.toString()} 
                            onValueChange={(value) => setLogSearch(prev => ({ ...prev, tail: Number(value) }))}
                        >
                            <SelectTrigger className="lg:w-full md:w-52 w-28  bg-gray-200/60">
                                <span>{logSearch.tail === 10000 ? 'All' : logSearch.tail}</span>
                            </SelectTrigger>
                            <SelectContent>
                                {[5, 10, 20, 500, 1000, 5000, 10000].map((value) => (
                                    <SelectItem key={value} value={value.toString()}>
                                        {value === 0 ? 'All' : value}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                    </div>
                </div>

                <div className="flex justify-between w-full mt-6">
                    <p className="text-gray-500 ml-1 w-1/6 whitespace-nowrap">{t('日志数据')}</p>
                    {/* 日志显示区域 */}
                    <div className="lg:w-full ml-6">
                        <Codemirror
                            ref={codemirrorRef}
                            value={logInfo}
                            editable={false}
                            width="85%"
                            height="950px"
                            theme="light"
                            autoFocus={true} // 加载自动聚焦
                            extensions={[javascript()]}             
                        />
                    </div>
                </div>
            </SheetContent>
        </Sheet>
    );
}



