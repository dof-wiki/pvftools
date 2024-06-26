package updater

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"pvftools/backend/common/log"
	"pvftools/backend/common/setting"
	"pvftools/backend/internal/progress"
	"strconv"
	"strings"
)

type Updater struct {
	curExecutable string
	curVersion    string
	newVersion    string
	newExecutable string
	updateMessage string
	newVersionUrl string
}

func (u *Updater) fetchVersion() {
	url := "https://api.github.com/repos/dof-wiki/pvftools/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		log.LogError("请求失败: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.LogError("请求失败: %v", resp.Status)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.LogError("读取失败: %v", err)
		return
	}
	jsonStr := string(body)
	u.newVersion = gjson.Get(jsonStr, "name").String()
	u.updateMessage = gjson.Get(jsonStr, "body").String()
	u.newVersionUrl = gjson.Get(jsonStr, "assets.0.browser_download_url").String()
}

func (u *Updater) CheckUpdate() bool {
	if u.newVersion == "" {
		u.fetchVersion()
	}
	return u.curVersion != u.newVersion
}

func (u *Updater) GetUpdateInfo() (string, string) {
	return u.newVersion, u.updateMessage
}

func (u *Updater) download() {
	resp, err := http.Get(u.newVersionUrl)
	if err != nil {
		log.LogError("请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		log.LogError("获取文件大小失败: %v", err)
		return
	}
	newExecutableName := fmt.Sprintf("pvftools_%s", strings.ReplaceAll(u.newVersion, ".", "_"))
	u.newExecutable = path.Join(setting.UpdaterDir(), newExecutableName)
	if _, err := os.Stat(u.newExecutable); err == nil {
		err = os.Remove(u.newExecutable)
		if err != nil {
			log.LogError("删除旧文件失败: %v", err)
			return
		}
	}

	out, err := os.Create(u.newExecutable)
	if err != nil {
		log.LogError("创建文件失败: %v", err)
		return
	}
	defer out.Close()
	pw := progress.NewProgressController("正在下载", size)
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	if err != nil {
		log.LogError("下载失败: %v", err)
		return
	}
	pw.End()
}

func (u *Updater) AutoUpdate() {
	if u.curVersion == u.newVersion {
		return
	}
	u.download()
	GenScript(u.curExecutable, u.newExecutable)
}

func (u *Updater) DoUpdate() {
	cmd := exec.Command("pvftools_updater.bat")
	err := cmd.Start()
	if err != nil {
		log.LogError("更新程序启动失败")
	}
	log.LogInfo("更新完成")
}

func NewUpdater() *Updater {
	path, _ := os.Executable()
	return &Updater{
		curExecutable: path,
		curVersion:    Version,
	}
}
