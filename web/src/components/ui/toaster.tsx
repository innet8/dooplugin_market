import { useToast } from "@/hooks/use-toast"
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  Toast,
  ToastClose,
  ToastDescription,
  ToastProvider,
  ToastTitle,
  ToastViewport,
} from "@/components/ui/toast"

export function SuccessToaster() {
  const { toasts } = useToast()

  return (
    <ToastProvider>
      {toasts.map(function ({ id, title, description, action, ...props }) {
        return (
          <Toast key={id} {...props}>
                <div className="grid gap-1">
                      <div className="flex items-center">
                        <Avatar className="size-6 mr-2">
                          <AvatarImage src="../../../public/安装成功.png"/>
                          <AvatarFallback />
                        </Avatar>
                        {title && <ToastTitle>{title}</ToastTitle>}
                      </div>
                      {description && (
                        <ToastDescription className="ml-7">{description}</ToastDescription>
                      )}
                
                {action}
            </div>
            
            <ToastClose />
          </Toast>
        )
      })}
      <ToastViewport />
    </ToastProvider>
  )
}

export function FalseToaster() {
  const { toasts } = useToast()

  return (
    <ToastProvider>
      {toasts.map(function ({ id, title, description, action, ...props }) {
        return (
          <Toast key={id} {...props}>
                <div className="grid gap-1">
                  <div className="flex items-center">
                      <Avatar className="size-8 mr-2">
                        <AvatarImage src="../../../public/安装失败.png"/>
                        <AvatarFallback />
                      </Avatar>
                        {title && <ToastTitle>{title}</ToastTitle>}
                    </div>
                  {description && (
                    <ToastDescription className="ml-10">{description}</ToastDescription>
                  )}
                
                {action}
            </div>
            
            <ToastClose />
          </Toast>
        )
      })}
      <ToastViewport />
    </ToastProvider>
  )
}
