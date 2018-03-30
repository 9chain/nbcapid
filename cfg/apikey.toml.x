[user]
configDir = "./userconfig"
maxUserFiles = 20

[session]
sessionDir = "/tmp/session_nbcsrv01"
maxageMin = 600
sessionKey = "session key"

[sdksrv]
username = "superuser"
apiKey = "api key"

[smtp]
host = "smtp.163.com"
serverAddr = "smtp.163.com:25"
user = "aquariusye@163.com"
password = "helloshiki"
salt = "salt"
timeoutMin = 120
pageForgetPassword = "/public/resetpassword.html"
confirmUrl = "http://localhost:8080/panel"
activeTitle = "注册确认邮件"
resetPasswordTitle = "重设密码确认邮件"