package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os_manage/azur_lane"
	"os_manage/log"
	"strings"
	"sync"
)

var (
	azureLaneBox   GroupBox
	azureLanePanel *AzureLanePanel

	biliSyncOnce      sync.Once
	azureLaneSyncOnce sync.Once

	rainbowColor = GradientBrush{
		Vertexes: []walk.GradientVertex{
			{X: 0, Y: 0, Color: walk.RGB(255, 255, 127)},
			{X: 1, Y: 0, Color: walk.RGB(127, 191, 255)},
			{X: 0.5, Y: 0.5, Color: walk.RGB(255, 255, 255)},
			{X: 1, Y: 1, Color: walk.RGB(127, 255, 127)},
			{X: 0, Y: 1, Color: walk.RGB(255, 127, 127)},
		},
		Triangles: []walk.GradientTriangle{
			{0, 1, 2},
			{1, 3, 2},
			{3, 4, 2},
			{4, 0, 2},
		},
	}
)

type AzureLanePanel struct {
	*walk.MainWindow
	HWnd win.HWND

	TypeBackRow      *walk.CheckBox // 后排主力
	TypeFrontRow     *walk.CheckBox // 前排先锋
	TypeBattleShip   *walk.CheckBox // 战列
	TypeCarrierShip  *walk.CheckBox // 航母
	TypeWeightPatrol *walk.CheckBox // 重巡
	TypeLightPatrol  *walk.CheckBox // 轻巡
	TypeDestroyer    *walk.CheckBox // 驱逐
	TypeSubmarine    *walk.CheckBox // 潜艇
	TypeSail         *walk.CheckBox // 风帆
	TypeOther        *walk.CheckBox // 重炮、维修、运输

	RarityNormal    *walk.CheckBox // 普通
	RarityRare      *walk.CheckBox // 稀有
	RarityElite     *walk.CheckBox // 精锐
	RaritySuperRare *walk.CheckBox // 超稀有
	RarityUltraRare *walk.CheckBox // 海上传奇
	RarityPriority  *walk.CheckBox // 最高方案
	RarityDecisive  *walk.CheckBox // 决战方案

	CampEagleUnion         *walk.CheckBox // 白鹰
	CampRoyalNavy          *walk.CheckBox // 皇家
	CampSakuraIslands      *walk.CheckBox // 重樱
	CampIronBlood          *walk.CheckBox // 铁血
	CampDragonEmpery       *walk.CheckBox // 东煌
	CampNorthernParliament *walk.CheckBox // 北方联合
	CampIrisTheLiberty     *walk.CheckBox // 自由鸢尾
	CampCuriaOfVichya      *walk.CheckBox // 维希教廷
	CampSardinianEmpire    *walk.CheckBox // 撒丁帝国
	CampOther              *walk.CheckBox // 其他
}

func (p *AzureLanePanel) Types() string {
	var types []string
	if p.TypeBackRow.Checked() {
		types = append(types, "后排主力")
	}
	if p.TypeFrontRow.Checked() {
		types = append(types, "前排先锋")
	}
	if p.TypeBattleShip.Checked() {
		types = append(types, "战")
	}
	if p.TypeCarrierShip.Checked() {
		types = append(types, "航")
	}
	if p.TypeWeightPatrol.Checked() {
		types = append(types, "重巡|超巡")
	}
	if p.TypeLightPatrol.Checked() {
		types = append(types, "轻巡")
	}
	if p.TypeDestroyer.Checked() {
		types = append(types, "驱逐")
	}
	if p.TypeSubmarine.Checked() {
		types = append(types, "潜")
	}
	if p.TypeSail.Checked() {
		types = append(types, "风帆")
	}
	if p.TypeOther.Checked() {
		types = append(types, "重炮|维修|运输")
	}

	if len(types) > 0 {
		return fmt.Sprintf(" type regexp '%s' ", strings.Join(types, "|"))
	}
	return ""
}

func (p *AzureLanePanel) Rarity() string {
	var rarities []string
	if p.RarityNormal.Checked() {
		rarities = append(rarities, "'普通'")
	}
	if p.RarityRare.Checked() {
		rarities = append(rarities, "'稀有'")
	}
	if p.RarityElite.Checked() {
		rarities = append(rarities, "'精锐'")
	}
	if p.RaritySuperRare.Checked() {
		rarities = append(rarities, "'超稀有', '最高方案'")
	}
	if p.RarityUltraRare.Checked() {
		rarities = append(rarities, "'海上传奇', '决战方案'")
	}

	if len(rarities) > 0 {
		return fmt.Sprintf("rarity in (%s)", strings.Join(rarities, ", "))
	}
	return ""
}

func (p *AzureLanePanel) Camp() string {
	var rarities []string
	if p.CampEagleUnion.Checked() {
		rarities = append(rarities, "'白鹰'")
	}
	if p.CampRoyalNavy.Checked() {
		rarities = append(rarities, "'皇家'")
	}
	if p.CampSakuraIslands.Checked() {
		rarities = append(rarities, "'重樱'")
	}
	if p.CampIronBlood.Checked() {
		rarities = append(rarities, "'铁血'")
	}
	if p.CampDragonEmpery.Checked() {
		rarities = append(rarities, "'东煌'")
	}
	if p.CampSardinianEmpire.Checked() {
		rarities = append(rarities, "'撒丁帝国'")
	}
	if p.CampNorthernParliament.Checked() {
		rarities = append(rarities, "'北方联合'")
	}
	if p.CampIrisTheLiberty.Checked() {
		rarities = append(rarities, "'自由鸢尾'")
	}
	if p.CampCuriaOfVichya.Checked() {
		rarities = append(rarities, "'维希教廷'")
	}
	if p.CampOther.Checked() {
		rarities = append(rarities, "") // todo
	}

	if len(rarities) > 0 {
		return fmt.Sprintf("camp in (%s)", strings.Join(rarities, ", "))
	}
	return ""
}

func (p *AzureLanePanel) Tech() {
	sql := `select * from ships`

	conditions := make([]string, 0)
	if p.Types() != "" {
		conditions = append(conditions, p.Types())
	}
	if p.Rarity() != "" {
		conditions = append(conditions, p.Rarity())
	}
	if p.Camp() != "" {
		conditions = append(conditions, p.Camp())
	}
	where := strings.Join(conditions, " and ")
	if where != "" {
		sql += " where " + where
	}
	sql += " order by tech_per_mind desc"
	log.Debug(sql)

	resp, err := azur_lane.LoadShips(sql)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("count", len(resp))
	for _, ship := range resp {
		log.Debug(ship.Name, ship.Type, ship.Camp, ship.Rarity, ship.TechPoint)
		break
	}
}

func getAzureLaneBox(alp *AzureLanePanel) GroupBox {
	azureLaneSyncOnce.Do(func() {
		typeLine1 := []Widget{
			CheckBox{
				AssignTo:  &alp.TypeBackRow,
				Text:      "后排主力",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeFrontRow,
				Text:      "前排先锋",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeBattleShip,
				Text:      "战列",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeCarrierShip,
				Text:      "航母",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeWeightPatrol,
				Text:      "重巡",
				OnClicked: alp.Tech,
			},
		}
		typeLine2 := []Widget{
			CheckBox{
				AssignTo:  &alp.TypeLightPatrol,
				Text:      "轻巡",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeDestroyer,
				Text:      "驱逐",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeSubmarine,
				Text:      "潜艇",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeSail,
				Text:      "风帆",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.TypeOther,
				Text:      "其他", // 重炮 维修 运输
				OnClicked: alp.Tech,
			},
		}
		typeGroup := GroupBox{
			Layout: VBox{}, Title: "类型",
			Children: []Widget{
				Composite{
					MaxSize: Size{Height: 24}, Layout: HBox{}, Children: typeLine1,
				},
				Composite{
					MaxSize: Size{Height: 24}, Layout: HBox{}, Children: typeLine2,
				},
			},
		}
		rarityGroup := GroupBox{
			Layout: HBox{}, Title: "稀有度",
			Children: []Widget{
				CheckBox{
					AssignTo:   &alp.RarityNormal,
					Background: SolidColorBrush{Color: 0xA8AEB5},
					Text:       "普通",
					OnClicked:  alp.Tech,
				}, CheckBox{
					AssignTo:   &alp.RarityRare,
					Background: SolidColorBrush{Color: 0xF3E67C},
					Text:       "稀有",
					OnClicked:  alp.Tech,
				}, CheckBox{
					AssignTo:   &alp.RarityElite,
					Background: SolidColorBrush{Color: 0xDDA0DD},
					Text:       "精锐",
					OnClicked:  alp.Tech,
				}, CheckBox{
					AssignTo:   &alp.RaritySuperRare,
					Background: SolidColorBrush{Color: 0x6BC7F7},
					Text:       "超稀有",
					OnClicked:  alp.Tech,
				}, CheckBox{
					AssignTo:   &alp.RarityUltraRare,
					Background: rainbowColor,
					Text:       "海上传奇",
					OnClicked:  alp.Tech,
				},
			},
		}
		campLine1 := []Widget{
			CheckBox{
				AssignTo:  &alp.CampEagleUnion,
				Text:      "白鹰",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampRoyalNavy,
				Text:      "皇家",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampSakuraIslands,
				Text:      "重樱",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampIronBlood,
				Text:      "铁血",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampDragonEmpery,
				Text:      "东煌",
				OnClicked: alp.Tech,
			},
		}
		campLine2 := []Widget{
			CheckBox{
				AssignTo:  &alp.CampSardinianEmpire,
				Text:      "撒丁帝国",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampNorthernParliament,
				Text:      "北方联合",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampIrisTheLiberty,
				Text:      "自由鸢尾",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampCuriaOfVichya,
				Text:      "维希教廷",
				OnClicked: alp.Tech,
			}, CheckBox{
				AssignTo:  &alp.CampOther,
				Text:      "其他",
				OnClicked: alp.Tech,
			},
		}
		campGroup := GroupBox{
			Layout: VBox{}, Title: "阵营",
			Children: []Widget{
				Composite{
					MaxSize: Size{Height: 24}, Layout: HBox{}, Children: campLine1,
				},
				Composite{
					MaxSize: Size{Height: 24}, Layout: HBox{}, Children: campLine2,
				},
			},
		}
		azureLaneBox = GroupBox{
			Layout: VBox{}, Title: "碧蓝航线科技点",
			Children: []Widget{typeGroup, rarityGroup, campGroup},
		}
	})

	return azureLaneBox
}

func OpenAzureLanePanel() {
	if azureLanePanel == nil {
		azureLanePanel = &AzureLanePanel{}
		err := MainWindow{
			AssignTo: &azureLanePanel.MainWindow, Size: Size{Width: 560, Height: 370},
			Layout: VBox{}, Title: "碧蓝航线科技",
			Children: []Widget{
				getAzureLaneBox(azureLanePanel),
				Composite{
					Layout: HBox{},
					Children: []Widget{
						PushButton{Text: "重置"},        // todo reset
						PushButton{Text: "更新版本船数据"},   // todo 更新版本船数据
						PushButton{Text: "过滤已120级的船"}, // todo 120级船记录过滤
					},
				},
				Composite{}, // todo 结果显示
			},
		}.Create()
		if err != nil {
			log.Error("create azure lane error:", err)
			return
		}
		azureLanePanel.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
			if reason == walk.CloseReasonUnknown {
				*canceled = true
				azureLanePanel.Hide()
			}
		})

		azureLanePanel.Run()
		return
	}

	azureLanePanel.Show()
}
