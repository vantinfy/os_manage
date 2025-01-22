package gui

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os_manage/azur_lane"
	"os_manage/log"
	"reflect"
	"sort"
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

	TechFill          *walk.CheckBox // 装填
	TechHit           *walk.CheckBox // 命中
	TechBombard       *walk.CheckBox // 炮击
	TechAviation      *walk.CheckBox // 航空
	TechMotorized     *walk.CheckBox // 机动
	TechThunder       *walk.CheckBox // 雷击
	TechAirDefence    *walk.CheckBox // 防空
	TechLasting       *walk.CheckBox // 耐久
	TechAntisubmarine *walk.CheckBox // 反潜

	tbView    *walk.TableView
	shipModel *ShipsModel
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
		rarities = append(rarities, "'飓风', '其他', 'META-???,其他'")
	}

	if len(rarities) > 0 {
		return fmt.Sprintf("camp in (%s)", strings.Join(rarities, ", "))
	}
	return ""
}

func (p *AzureLanePanel) Teches() string {
	var teches []string
	if p.TechFill.Checked() {
		teches = append(teches, "装填")
	}
	if p.TechHit.Checked() {
		teches = append(teches, "命中")
	}
	if p.TechBombard.Checked() {
		teches = append(teches, "炮击")
	}
	if p.TechAviation.Checked() {
		teches = append(teches, "航空")
	}
	if p.TechMotorized.Checked() {
		teches = append(teches, "机动")
	}
	if p.TechThunder.Checked() {
		teches = append(teches, "雷击")
	}
	if p.TechAirDefence.Checked() {
		teches = append(teches, "防空")
	}
	if p.TechLasting.Checked() {
		teches = append(teches, "耐久")
	}
	if p.TechAntisubmarine.Checked() {
		teches = append(teches, "反潜")
	}

	if len(teches) > 0 {
		return fmt.Sprintf(" tech_point regexp '%s' ", strings.Join(teches, "|"))
	}
	return ""
}

func (p *AzureLanePanel) TechQuery() {
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
	if p.Teches() != "" {
		conditions = append(conditions, p.Teches())
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
	if p.shipModel != nil {
		shipList := make([]shipItem, 0)
		for _, ship := range resp {
			shipList = append(shipList, shipItem{Ship: ship})
		}

		p.shipModel.items = shipList
		p.shipModel.PublishRowsReset()

		_ = p.shipModel.Sort(p.shipModel.sortColumn, p.shipModel.sortOrder)
	}
}

func (p *AzureLanePanel) ResetSelect() {
	pType := reflect.TypeOf(*p)
	vType := reflect.ValueOf(*p)
	checkBox := &walk.CheckBox{}
	// 所有复选框置否
	for i := 0; i < pType.NumField(); i++ {
		if pType.Field(i).Type == reflect.TypeOf(checkBox) {
			if !vType.Field(i).IsZero() {
				checkMethod := vType.Field(i).MethodByName("SetChecked")
				if checkMethod.IsValid() {
					checkMethod.Call([]reflect.Value{reflect.ValueOf(false)})
				}
			}
		}
	}

	// 清空所有船
	p.shipModel.items = make([]shipItem, 0)
	p.shipModel.PublishRowsReset()
}

func (p *AzureLanePanel) UpdateShips() {
	err := azur_lane.GetAllShips()
	if err != nil {
		log.Error(err)
	}
	log.Info("舰船数据更新完成")
}

func getAzureLaneBox(alp *AzureLanePanel) GroupBox {
	azureLaneSyncOnce.Do(func() {
		typeLine1 := []Widget{
			CheckBox{
				AssignTo:  &alp.TypeBackRow,
				Text:      "后排主力",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeFrontRow,
				Text:      "前排先锋",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeBattleShip,
				Text:      "战列",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeCarrierShip,
				Text:      "航母",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeWeightPatrol,
				Text:      "重巡",
				OnClicked: alp.TechQuery,
			},
		}
		typeLine2 := []Widget{
			CheckBox{
				AssignTo:  &alp.TypeLightPatrol,
				Text:      "轻巡",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeDestroyer,
				Text:      "驱逐",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeSubmarine,
				Text:      "潜艇",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeSail,
				Text:      "风帆",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TypeOther,
				Text:      "其他", // 重炮 维修 运输
				OnClicked: alp.TechQuery,
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
					OnClicked:  alp.TechQuery,
				}, CheckBox{
					AssignTo:   &alp.RarityRare,
					Background: SolidColorBrush{Color: 0xF3E67C},
					Text:       "稀有",
					OnClicked:  alp.TechQuery,
				}, CheckBox{
					AssignTo:   &alp.RarityElite,
					Background: SolidColorBrush{Color: 0xDDA0DD},
					Text:       "精锐",
					OnClicked:  alp.TechQuery,
				}, CheckBox{
					AssignTo:   &alp.RaritySuperRare,
					Background: SolidColorBrush{Color: 0x6BC7F7},
					Text:       "超稀有",
					OnClicked:  alp.TechQuery,
				}, CheckBox{
					AssignTo:   &alp.RarityUltraRare,
					Background: rainbowColor,
					Text:       "海上传奇",
					OnClicked:  alp.TechQuery,
				},
			},
		}
		campLine1 := []Widget{
			CheckBox{
				AssignTo:  &alp.CampEagleUnion,
				Text:      "白鹰",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampRoyalNavy,
				Text:      "皇家",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampSakuraIslands,
				Text:      "重樱",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampIronBlood,
				Text:      "铁血",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampDragonEmpery,
				Text:      "东煌",
				OnClicked: alp.TechQuery,
			},
		}
		campLine2 := []Widget{
			CheckBox{
				AssignTo:  &alp.CampSardinianEmpire,
				Text:      "撒丁帝国",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampNorthernParliament,
				Text:      "北方联合",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampIrisTheLiberty,
				Text:      "自由鸢尾",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampCuriaOfVichya,
				Text:      "维希教廷",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.CampOther,
				Text:      "其他",
				OnClicked: alp.TechQuery,
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
		techLine1 := []Widget{
			CheckBox{
				AssignTo:  &alp.TechFill,
				Text:      "装填",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechHit,
				Text:      "命中",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechBombard,
				Text:      "炮击",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechAviation,
				Text:      "航空",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechMotorized,
				Text:      "机动",
				OnClicked: alp.TechQuery,
			},
		}
		techLine2 := []Widget{
			CheckBox{
				AssignTo:  &alp.TechThunder,
				Text:      "雷击",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechAirDefence,
				Text:      "防空",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechLasting,
				Text:      "耐久",
				OnClicked: alp.TechQuery,
			}, CheckBox{
				AssignTo:  &alp.TechAntisubmarine,
				Text:      "反潜",
				OnClicked: alp.TechQuery,
			},
		}
		techGroup := GroupBox{
			Layout: VBox{}, Title: "科技点",
			Children: []Widget{
				Composite{
					MaxSize: Size{Height: 24}, Layout: HBox{}, Children: techLine1,
				},
				Composite{
					MaxSize: Size{Height: 24}, Layout: HBox{}, Children: techLine2,
				},
			},
		}
		azureLaneBox = GroupBox{
			Layout: VBox{}, Title: "碧蓝航线科技点",
			Children: []Widget{typeGroup, rarityGroup, campGroup, techGroup},
		}
	})

	return azureLaneBox
}

func OpenAzureLanePanel() {
	if azureLanePanel == nil {
		azureLanePanel = &AzureLanePanel{
			shipModel: &ShipsModel{},
		}
		err := MainWindow{
			AssignTo: &azureLanePanel.MainWindow, Size: Size{Width: 720, Height: 560},
			Layout: VBox{}, Title: "碧蓝航线科技",
			Children: []Widget{
				getAzureLaneBox(azureLanePanel),
				Composite{
					Layout: HBox{},
					Children: []Widget{
						PushButton{Text: "重置", OnClicked: azureLanePanel.ResetSelect},
						PushButton{Text: "更新版本船数据", OnClicked: azureLanePanel.UpdateShips},
						PushButton{Text: "过滤已120级的船"}, // todo 120级船记录过滤
					}, // todo 逆转日志文本区 快捷键 便签
				},
				TableView{
					AssignTo:         &azureLanePanel.tbView,
					AlternatingRowBG: true,
					ColumnsOrderable: true,
					CheckBoxes:       true,
					MultiSelection:   true,
					Columns: []TableViewColumn{
						{Name: "Name", DataMember: "名称"},
						//{Name: "Avatar", DataMember: "头像"},
						{Name: "Type", DataMember: "类型"},
						{Name: "TechPoint", DataMember: "科技点"},
						{Name: "Camp", DataMember: "阵营"},
						{Name: "Rarity", DataMember: "稀有度"},
						{Name: "MindCost", DataMember: "120级消耗心智Ⅰ"},
						{Name: "TechPerMind", DataMember: "科技点/心智Ⅰ"},
					},
					Model: azureLanePanel.shipModel,
					OnSelectedIndexesChanged: func() {
						fmt.Printf("SelectedIndexes: %v\n", azureLanePanel.tbView.SelectedIndexes())
					},
				}, // 结果显示
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

type ShipsModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items []shipItem
}

type shipItem struct {
	azur_lane.Ship
	checked bool
}

func (m *ShipsModel) Items() any {
	return m.items
}

func (m *ShipsModel) RowCount() int {
	return len(m.items)
}

func (m *ShipsModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Name

	//case 1:
	//	return item.Avatar

	case 1:
		return item.Type

	case 2:
		return item.TechPoint

	case 3:
		return item.Camp

	case 4:
		return item.Rarity

	case 5:
		return item.MindCost

	case 6:
		return fmt.Sprintf("%.6f", item.TechPerMind)
	}

	panic("unexpected col")
}

func (m *ShipsModel) Checked(row int) bool {
	return m.items[row].checked
}

func (m *ShipsModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

func (m *ShipsModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order

	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]

		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}

			return !ls
		}

		switch m.sortColumn {
		case 0:
			return c(a.Name < b.Name)

		//case 1:
		//	return c(a.Avatar < b.Avatar)

		case 1:
			if len(a.Type) == len(b.Type) {
				for k := 0; k < len(a.Type); k++ {
					if a.Type[k] != b.Type[k] {
						return c(a.Type[k] < b.Type[k])
					}
				}
			}
			return false

		case 2:
			if len(a.TechPoint) == len(b.TechPoint) {
				for k := 0; k < len(a.TechPoint); k++ {
					if a.TechPoint[k] != b.TechPoint[k] {
						return c(a.TechPoint[k] < b.TechPoint[k])
					}
				}
			}
			return false

		case 3:
			return c(a.Camp < b.Camp)

		case 4:
			return c(a.Rarity < b.Rarity)

		case 5:
			return c(a.MindCost < b.MindCost)

		case 6:
			return c(fmt.Sprintf("%.6f", a.TechPerMind) < fmt.Sprintf("%.6f", b.TechPerMind))

		}

		panic("unreachable")
	})

	return m.SorterBase.Sort(col, order)
}
