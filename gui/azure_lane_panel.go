package gui

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
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

	TypeBattleShip   *walk.CheckBox // 战列
	TypeCarrierShip  *walk.CheckBox // 航母
	TypeWeightPatrol *walk.CheckBox // 重巡
	TypeLightPatrol  *walk.CheckBox // 轻巡
	TypeDestroyer    *walk.CheckBox // 驱逐
	TypeSubmarine    *walk.CheckBox // 潜艇
	TypeSail         *walk.CheckBox // 风帆

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

func getAzureLaneBox(alp *AzureLanePanel) GroupBox {
	azureLaneSyncOnce.Do(func() {
		typeGroup := GroupBox{
			Layout: HBox{}, Title: "类型",
			Children: []Widget{
				CheckBox{
					Text: "战列",
				}, CheckBox{
					Text: "航母",
				}, CheckBox{
					Text: "重巡",
				}, CheckBox{
					Text: "轻巡",
				}, CheckBox{
					Text: "驱逐",
				}, CheckBox{
					Text: "潜艇",
				},
			},
		}
		rarityGroup := GroupBox{
			Layout: HBox{}, Title: "稀有度",
			Children: []Widget{
				CheckBox{
					Background: SolidColorBrush{Color: 0xA8AEB5},
					Text:       "普通",
				}, CheckBox{
					Background: SolidColorBrush{Color: 0xF3E67C},
					Text:       "稀有",
				}, CheckBox{
					Background: SolidColorBrush{Color: 0xDDA0DD},
					Text:       "精锐",
				}, CheckBox{
					Background: SolidColorBrush{Color: 0x6BC7F7},
					Text:       "超稀有",
				}, CheckBox{
					Background: rainbowColor,
					Text:       "海上传奇",
				},
			},
		}
		campLine1 := []Widget{
			CheckBox{
				Text: "白鹰",
			}, CheckBox{
				Text: "皇家",
			}, CheckBox{
				Text: "重樱",
			}, CheckBox{
				Text: "铁血",
			}, CheckBox{
				Text: "东煌",
			},
		}
		campLine2 := []Widget{
			CheckBox{
				Text: "撒丁帝国",
			}, CheckBox{
				Text: "北方联合",
			}, CheckBox{
				Text: "自由鸢尾",
			}, CheckBox{
				Text: "维希教廷",
			}, CheckBox{
				Text: "其他",
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
				PushButton{Text: "重置"}, // todo reset
				Composite{},            // todo 结果
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
