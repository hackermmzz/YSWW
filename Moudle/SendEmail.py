import smtplib
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.header import Header
from  email.utils import formataddr
from SockConn import SocketConn
import json
import argparse
import SockConn
def send_email(sender_email, sender_password, receiver_email, subject, content, sender_name=None):
    # 创建一个带附件的邮件实例
    message = MIMEMultipart()

    # 设置发件人字段，使用formataddr函数确保格式正确
    if sender_name:
        # 对发件人名称进行编码，确保中文正确显示
        sender_display = formataddr((str(Header(sender_name, 'utf-8')), sender_email))
    else:
        # 如果没有提供显示名称，则只使用邮箱地址
        sender_display = sender_email
        
    message['From'] = sender_display
    message['To'] = receiver_email
    message['Subject'] = Header(subject, 'utf-8')

    # 添加邮件正文
    message.attach(MIMEText(content, 'plain', 'utf-8'))

    try:
        # 创建 SMTP 对象并连接到 SMTP 服务器
        smtp_obj = smtplib.SMTP_SSL('smtp.qq.com', 465)

        # 登录发件人邮箱
        smtp_obj.login(sender_email, sender_password)

        # 发送邮件
        smtp_obj.sendmail(sender_email, receiver_email, message.as_string())

        print("邮件发送成功")
        return True

    except smtplib.SMTPException as e:
        print(f"邮件发送失败，错误信息: {e}")
        return False
    finally:
        # 关闭 SMTP 连接
        if 'smtp_obj' in locals():
            smtp_obj.quit()

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='SendEmail')
    parser.add_argument('--port',required=True,type=int)
    args = parser.parse_args()
    #连接后端
    conn=SocketConn(args.port)
    conn.write(SockConn.UniverseCookie)
    # 请替换为你的邮箱信息
    sender_email = "2049983474@qq.com"  # 发件人邮箱
    sender_password = "eumlxuysimsrdjbf"  # 发件人邮箱授权码
    sender_name="一生万物官网"
    #与主机建立联系
    while True:
        receiver_email = conn.read()#账号邮箱
        subject=conn.read() #邮箱主题
        content=conn.read() #邮箱内容
        status=send_email(sender_email, sender_password, receiver_email, subject, content,sender_name)
        
        #返回邮件信息给服务器
        msg={"status":status,"email":receiver_email}
        js=json.dumps(msg)
        conn.write(js)