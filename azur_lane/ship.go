package azur_lane

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"net/url"
	"os"
	"os_manage/database"
	"os_manage/log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Ship struct {
	Name           string   `json:"name"`             // 船名
	Avatar         string   `json:"avatar"`           // 头像
	Clothing       []string `json:"clothing"`         // 服装
	Type           []string `json:"type"`             // 前后排 种类
	TechPointCount int      `json:"tech_point_count"` // 阵营科技总和
	TechPointCamp  []int    `json:"tech_point_camp"`  // 阵营科技点加成 gain max lv120
	TechPoint      []string `json:"tech_point"`       // 全局加成 gain - lv120
	Rarity         string   `json:"rarity"`           // 稀有度
	MindCost       int      `json:"mind_cost"`        // 120级心智消耗 所需金币*100即可
	Camp           string   `json:"camp"`             // 阵营
	ConstructTime  string   `json:"construct_time"`   // 建造时间
	InstallDate    string   `json:"install_date"`     // 实装日期
	TransformDate  string   `json:"transform_date"`   // 改造日期
	OrdinaryDrop   []string `json:"ordinary_drop"`    // 获取方式-普通掉落点
	FileDrop       []string `json:"file_drop"`        // 获取方式-档案掉落点
	ActivityDrop   []string `json:"activity_drop"`    // 获取方式-活动掉落点
	Cute           string   `json:"cute"`             // 默认Q版立绘
	TechPerMind    float64  `json:"tech_per_mind"`    // 每一点心智一提升多少科技点 越高越好
	Has120         bool     `json:"has_120"`          // 该船已经练到120级
}

const (
	shipInfoPagePrefix = `https://wiki.biligame.com/blhx/`
	shipDetailSavePath = "ships_detail/"
)

func GetShipDetailInfo(ship *Ship) error {
	if len(ship.Name) <= 0 {
		return fmt.Errorf("ship name is empty")
	}

	// https://wiki.biligame.com/blhx/%E7%89%B9%E8%A3%85%E5%9E%8B%E5%B8%83%E9%87%8CMKIII
	req, _ := http.NewRequest(http.MethodGet, shipInfoPagePrefix+url.QueryEscape(ship.Name), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("get ship detail[%s] response: %s", ship.Name, resp.Status))

	AnalyseTechPoint(string(respBytes), ship)
	AnalyseClothing(string(respBytes), ship)
	AnalyseObtainSource(string(respBytes), ship)

	if _, err = os.Stat(shipDetailSavePath); os.IsNotExist(err) {
		_ = os.MkdirAll(shipDetailSavePath, 0644)
	}
	_ = os.WriteFile(shipDetailSavePath+ship.Name, respBytes, 0644)

	return nil
}

func GetShipDetailInfoLocal(ship *Ship) error {
	if len(ship.Name) <= 0 {
		return fmt.Errorf("ship name is empty")
	}

	shipDetails, err := os.ReadFile(shipDetailSavePath + ship.Name)
	if err != nil {
		return err
	}

	AnalyseTechPoint(string(shipDetails), ship)
	AnalyseClothing(string(shipDetails), ship)
	AnalyseObtainSource(string(shipDetails), ship)

	return nil
}

func AnalyseTechPoint(originRespHTML string, ship *Ship) {
	gainRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b>获得</b>\s*</td>\s*<td>\s*<table[^>]*>\s*<tbody><tr>\s*<td[^>]*><img[^>]*>\s*</td>\s*<td>\+\s*</td>\s*<td[^>]*>(.+?)\s*</td>\s*</tr></tbody></table>\s*</td>\s*<td[^>]*>合计<br\s*/>(.+?)\s*</td>\s*<td>(.+?)\s*</td>\s*</tr>`)
	gainMatch := gainRegexp.FindStringSubmatch(originRespHTML)
	if len(gainMatch) > 3 {
		ship.TechPointCount, _ = strconv.Atoi(gainMatch[2])

		gainTech, _ := strconv.Atoi(gainMatch[1])
		ship.TechPointCamp = append(ship.TechPointCamp, gainTech)

		ship.TechPoint = append(ship.TechPoint, gainMatch[3])
	}

	maxRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b>满星</b>\s*</td>\s*<td>\s*<table[^>]*>\s*<tbody><tr>\s*<td[^>]*><img[^>]*>\s*</td>\s*<td[^>]*>\+\s*</td>\s*<td[^>]*>(.+?)\s*</td>\s*</tr></tbody></table>\s*</td>\s*<td>(.+?)\s*</td>\s*</tr>`)
	maxMatch := maxRegexp.FindStringSubmatch(originRespHTML)
	if len(maxMatch) > 2 {
		maxTech, _ := strconv.Atoi(maxMatch[1])
		ship.TechPointCamp = append(ship.TechPointCamp, maxTech)

		ship.TechPoint = append(ship.TechPoint, maxMatch[2])
	}

	lv120Regexp := regexp.MustCompile(`(?s)<tr>\s*<td><b>Lv.120</b>\s*</td>\s*<td>\s*<table[^>]*>\s*<tbody><tr>\s*<td[^>]*><img[^>]*>\s*</td>\s*<td[^>]*>\+\s*</td>\s*<td[^>]*>(.+?)\s*</td>\s*</tr></tbody></table>\s*</td>\s*<td>(.+?)\s*</td>\s*</tr>`)
	lv120Match := lv120Regexp.FindStringSubmatch(originRespHTML)
	if len(lv120Match) > 2 {
		lv120Tech, _ := strconv.Atoi(lv120Match[1])
		ship.TechPointCamp = append(ship.TechPointCamp, lv120Tech)

		ship.TechPoint = append(ship.TechPoint, lv120Match[2])
	}
}

func AnalyseClothing(originRespHTML string, ship *Ship) {
	skinRegexp := regexp.MustCompile(`<div class="tab_con( active)?"><img alt=".+?.jpg" src="(.+?)" .+?/></div>`)
	cuteRegexp := regexp.MustCompile(`Q版立绘.png" src="(.+?)"`)

	skinMatch := skinRegexp.FindAllStringSubmatch(originRespHTML, -1)
	if len(skinMatch) > 0 {
		for _, match := range skinMatch {
			if len(match) > 2 {
				ship.Clothing = append(ship.Clothing, match[2])
			}
		}
	}

	cuteMatch := cuteRegexp.FindStringSubmatch(originRespHTML)
	// fmt.Println(cuteMatch, len(cuteMatch))
	if len(cuteMatch) > 1 {
		ship.Cute = cuteMatch[1]
	}
}

func AnalyseObtainSource(originRespHTML string, ship *Ship) {
	installDateRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">实装</span><span style="display:inline-block">日期</span></b>\s*</td>\s*<td colspan="4">(.+?)\s*</td></tr>`)
	transformDateRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">改造</span><span style="display:inline-block">日期</span></b>\s*</td>\s*<td colspan="4">(.+?)\s*</td></tr>`)
	constructTimeRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">建造</span><span style="display:inline-block">时间</span></b>\s*</td>\s*<td colspan="4">(.+?)\s*</td></tr>`)
	constructTimeRegexp2 := regexp.MustCompile(`建造时间">(.+)</a>`)

	// ordinaryDrop := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">普通</span><span style="display:inline-block">掉落点</span></b>\s*</td>\s*<td colspan="4"><div style="display:block;max-height:100px;overflow-y:auto">(<a href="/blhx/.+?" title=".+">(.+?)</a>、?)+?</div>\s*</td></tr>`)
	ordinaryDropRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">普通</span><span style="display:inline-block">掉落点</span></b>\s*</td>\s*<td colspan="4">(.+?)\s*</td></tr>`)
	fileDropRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">档案</span><span style="display:inline-block">掉落点</span></b>\s*</td>\s*<td colspan="4">(.+?)\s*</td></tr>`)
	activityDropRegexp := regexp.MustCompile(`(?s)<tr>\s*<td><b><span style="display:inline-block">活动</span><span style="display:inline-block">掉落点</span></b>\s*</td>\s*<td colspan="4">(.+?)\s*</td></tr>`)
	dropRegexp := regexp.MustCompile(`<a [^>]*>(.*?)</a>`)

	installMatch := installDateRegexp.FindStringSubmatch(originRespHTML)
	if len(installMatch) > 1 {
		ship.InstallDate = installMatch[1]
	}
	transformMatch := transformDateRegexp.FindStringSubmatch(originRespHTML)
	if len(transformMatch) > 1 {
		ship.TransformDate = transformMatch[1]
	}
	constructMatch := constructTimeRegexp.FindStringSubmatch(originRespHTML)
	if len(constructMatch) > 1 {
		if strings.Contains(constructMatch[1], "\u003c") {
			constructMatch := constructTimeRegexp2.FindStringSubmatch(constructMatch[1])
			if len(constructMatch) > 1 {
				ship.ConstructTime = constructMatch[1]
			}
		} else {
			ship.ConstructTime = constructMatch[1]
		}
	}

	ordinaryMatch := ordinaryDropRegexp.FindStringSubmatch(originRespHTML)
	if len(ordinaryMatch) > 1 {
		drops := dropRegexp.FindAllStringSubmatch(ordinaryMatch[1], -1)
		for _, v := range drops {
			if len(v) > 1 {
				ship.OrdinaryDrop = append(ship.OrdinaryDrop, v[1])
			}
		}
	}
	fileMatch := fileDropRegexp.FindStringSubmatch(originRespHTML)
	if len(fileMatch) > 1 {
		drops := dropRegexp.FindAllStringSubmatch(fileMatch[1], -1)
		for _, v := range drops {
			if len(v) > 1 {
				ship.FileDrop = append(ship.FileDrop, v[1])
			}
		}
	}
	activityMatch := activityDropRegexp.FindStringSubmatch(originRespHTML)
	if len(activityMatch) > 1 {
		drops := dropRegexp.FindAllStringSubmatch(activityMatch[1], -1)
		for _, v := range drops {
			if len(v) > 1 {
				ship.ActivityDrop = append(ship.ActivityDrop, v[1])
			}
		}
	}
}

func GetAllShips() error {
	respBytes, err := os.ReadFile("wiki.html")
	if err != nil {
		req, _ := http.NewRequest(http.MethodGet, `https://wiki.biligame.com/blhx/%E8%88%B0%E8%88%B9%E5%9B%BE%E9%89%B4`, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("获取wiki首页失败", err)
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("读取wiki首页响应数据失败", err)
			return err
		}
		err = os.WriteFile("wiki.html", body, 0644)
		if err != nil {
			log.Error("保存wiki首页文件失败", err)
			return err
		}
		respBytes = body
	}
	resp := string(respBytes)

	reg := regexp.MustCompile(`<div class="jntj-1 divsort" data-param0="0" .+?</a></span></div>`)
	list := reg.FindAllString(resp, -1)

	outList := make([]string, 0)
	for _, v := range list {
		if strings.Contains(v, "改") || strings.Contains(v, "联动") {
			continue
		}
		outList = append(outList, v)
	}
	log.Info("版本非联动舰船总数", len(outList))

	avatarReg := regexp.MustCompile(`https://patchwiki.biligame.com/images/blhx/thumb/.+/60px-.+?头像.jpg`)
	nameReg := regexp.MustCompile(`-(.+)头像`)
	typeReg := regexp.MustCompile(`data-param1="(.+),,(.+?)"`)
	rarityReg := regexp.MustCompile(`data-param2="(.+?)"`)
	campReg := regexp.MustCompile(`data-param3="(.+?)"`)

	//wg := sync.WaitGroup{}
	// 创建或加载缓存实例
	newCache, err := NewCache(CacheName)
	if err != nil {
		fmt.Println("Error loading cache:", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	cnt := 0 // 已获取数据的船数量
	go func() {
		ticker := time.NewTicker(7 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("舰船数据获取完毕")
				return
			case <-ticker.C:
				log.Info(fmt.Sprintf("获取舰船数据中[%d/%d]...请耐心等待", cnt, len(outList)))
			}
		}
	}()

	allShips := map[string]Ship{}
	for _, v := range outList {
		cnt++

		ship := &Ship{}
		htmlOrigin, err := url.QueryUnescape(v)
		if err != nil {
			log.Error(err)
			continue
		}

		ship.Avatar = avatarReg.FindString(htmlOrigin)
		nameMatch := nameReg.FindStringSubmatch(ship.Avatar)
		if len(nameMatch) > 1 {
			ship.Name = nameMatch[1]
		}
		// 缓存中已经存在
		if _, ok := newCache.Get(ship.Name); ok {
			ship, ok := newCache.Get(ship.Name)
			if ok {
				allShips[ship.Name] = *ship
			}
			continue
		}

		typeMatch := typeReg.FindStringSubmatch(htmlOrigin)
		if len(typeMatch) > 2 {
			ship.Type = []string{typeMatch[1], typeMatch[2]}
		}
		rarityMatch := rarityReg.FindStringSubmatch(htmlOrigin)
		if len(rarityMatch) > 1 {
			ship.Rarity = rarityMatch[1]
			switch ship.Rarity {
			case "普通":
				ship.MindCost = 660
			case "稀有":
				ship.MindCost = 880
			case "精锐":
				ship.MindCost = 1320
			case "超稀有", "最高方案":
				ship.MindCost = 2200
			case "海上传奇", "决战方案":
				ship.MindCost = 3300
			}
		}
		campMatch := campReg.FindStringSubmatch(htmlOrigin)
		if len(campMatch) > 1 {
			ship.Camp = campMatch[1]
		}

		// 不使用协程的情况下681艘船需要2分10秒左右 但是开协程会被ban ip
		//wg.Add(1)
		//go func(ship *Ship) {
		//	defer wg.Done()

		// 本地有已经下载好的单船详情
		if _, err = os.Stat(shipDetailSavePath + ship.Name); err == nil {
			err = GetShipDetailInfoLocal(ship)
			if err != nil {
				fmt.Printf("get ship[%s] detail by local file error: %v\n", ship.Name, err)
			}
		} else {
			err = GetShipDetailInfo(ship)
			if err != nil {
				fmt.Printf("get ship[%s] detail failed: %v\n", ship.Name, err)
			}
		}
		lv120 := 0
		if len(ship.TechPoint) == 3 {
			points := strings.Split(ship.TechPoint[2], "+")
			if len(points) > 1 {
				lv120, _ = strconv.Atoi(points[1])
			}
		}
		decLv120 := decimal.NewFromFloat(float64(lv120))
		decCost := decimal.NewFromFloat(float64(ship.MindCost))
		ship.TechPerMind, _ = decLv120.Div(decCost).Float64()
		//}(ship)

		allShips[ship.Name] = *ship
		newCache.Set(ship.Name, ship)
	}
	cancel()
	//wg.Wait()
	log.Info("get all ships: real ships count", len(newCache.Data))

	// 保存缓存到文件
	if err = newCache.Save(); err != nil {
		log.Error("Error saving cache:", err)
	}

	// 保存数据库
	db, err := database.GetDB(LocalDBName)
	if err != nil {
		log.Error("get all ships: get db failed", err)
		return err
	}
	defer db.Close()
	err = InsertShipData(db, allShips)
	if err != nil {
		log.Error("get all ships: insert ships to db failed", err)
		return err
	}

	return nil
}
