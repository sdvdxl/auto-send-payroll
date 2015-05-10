package main

import (
	"fmt"
	"github.com/sdvdxl/auto-send-payroll/config"
	"github.com/sdvdxl/auto-send-payroll/mail"
	"github.com/axgle/pinyin"
	"github.com/tealeg/xlsx"
	"strings"
	"os"
)

var (
	cfg   *config.Config
	works chan bool
	count int
)

func init() {
	var err error
	works = make(chan bool, 10)

	configFile :="config/config.yaml"
	fmt.Println("读取配置文件",configFile)
	cfg, err = config.ReadConfig(configFile)
	if cfg.Smtp_Server=="" {
		fmt.Println("smtp server 没有配置")
		os.Exit(1)
	}

	if cfg.Port==0 {
		fmt.Println("smtp server port 没有配置")
		os.Exit(1)
	}

	if cfg.Sender_Email =="" {
		fmt.Println("发件人 没有配置")
		os.Exit(1)
	}

	if cfg.Sender_Password=="" {
		fmt.Println("发件人名称 没有配置")
		os.Exit(1)
	}

	printErrorExit(err, "打开配置文件出错")
}

func main() {
	var (
		startRow, endRow int
		found            bool
		title            = [...]string{"部门", "姓名", "基本工资", "岗位工资", "通讯费",
			"病假", "事假", "年休", "缺勤", "扣款", "应付工资", "社保基数",
			"公积金基数", "养老", "失业", "医疗", "公积金", "个人合计",
			"免税金额", "补扣款", "应发工资", "应税金额",
			"代扣个税", "核发总额",
		}
	)

	excelFileName := cfg.Execl_Path
	xlFile, err := xlsx.OpenFile(excelFileName)
	fmt.Println("读取Excel", excelFileName)
	printErrorExit(err, "打开Excel文件出错")

	for x, row := range xlFile.Sheets[0].Rows {
		if found {
			break
		}

		for _, cell := range row.Cells {
			if "基本信息" == cell.Value {
				startRow = x
				break
			}

			if "合计" == cell.Value {
				endRow = x
				found = true
				break
			}
		}
	}

	rows := xlFile.Sheets[0].Rows[startRow+2 : endRow]
	count = len(rows)
	fmt.Println(fmt.Sprintf("一共%v个人",count))
	if count==0 {
		fmt.Println("没有可以发送的人员")
		os.Exit(0)
	}
	for i, row := range rows {
		payRollInfo := `<!Doctype html><html xmlns=http://www.w3.org/1999/xhtml><meta http-equiv=Content-Type content="text/html;charset=utf-8"><body><table cellspacing='0' border='1'><thead>`

		var name string
		for _, v := range title {
			if v == "" {
				v = "0"
			}

			payRollInfo += "<td nowrap>" + v + "</td>"
		}
		payRollInfo += "</th></thead><tbody><tr>"

		for j, cell := range row.Cells {
			if j == 0 {
				continue
			}

			if j == 2 {
				name = cell.Value
			}

			payRollInfo += "<td nowrap>" + cell.Value + "</td>"
		}
		payRollInfo += "</tr></tbody></table></body></html>"

		//将名字变换成拼音，
		nameChars := []rune(name)

		lastName := pinyin.Convert(string(nameChars[1:]))
		firstName := pinyin.Convert(string(nameChars[:1]))

		email := strings.ToLower(lastName + "." + firstName + "@fraudmetrix.cn")

		works <- true

		fmt.Println(fmt.Sprintf("正在给%s(%s)发送邮件", name, email))
		mail := mail.Mail{SmtpServer: cfg.Smtp_Server, Port: cfg.Port,
			SenderEmail: cfg.Sender_Email, SenderName: cfg.Sender_Name,
			SenderPassword: cfg.Sender_Password, ToEmail: email,
			Subject: cfg.Subject, Message: payRollInfo,
		}
		works<-true
		go func(num int) {
			mail.SendEmail()
			if num+1 == count {
				close(works)
			}
		}(i)
	}

	for _ = range works {
	}

	fmt.Println("发送完毕")
}

func printErrorExit(err error, msg string) {
	if err!=nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}