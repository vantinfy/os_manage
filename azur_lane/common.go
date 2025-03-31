package azur_lane

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "modernc.org/sqlite"
	"os_manage/database"
	"os_manage/log"
)

const (
	LocalDBName = "./ships.db"
	TableShip   = "ships"
	CacheName   = ".ships_cache"
)

func init() {
	// 连接 SQLite 数据库（如果文件不存在会自动创建）
	db, err := database.GetDB(LocalDBName)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	database.RegisterRegexp()

	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
    id INTEGER PRIMARY KEY AUTOINCREMENT,        -- 唯一ID，自增
    name TEXT NOT NULL constraint ships_pk unique,  -- 船名 唯一
    avatar TEXT,                                  -- 头像
    clothing TEXT,                                -- 服装，存储为JSON字符串
    type TEXT,                                    -- 前后排种类，存储为JSON字符串
    tech_point_count INTEGER DEFAULT 0,           -- 阵营科技总和
    tech_point_camp TEXT,                         -- 阵营科技点加成，存储为JSON字符串
    tech_point TEXT,                              -- 全局加成，存储为JSON字符串
    rarity TEXT,                                  -- 稀有度
    mind_cost INTEGER,                            -- 120级心智消耗
    camp TEXT,                                    -- 阵营
    construct_time TEXT,                          -- 建造时间（存储为字符串格式）
    install_date TEXT,                            -- 实装日期（存储为字符串格式）
    transform_date TEXT,                          -- 改造日期（存储为字符串格式）
    ordinary_drop TEXT,                           -- 普通掉落点，存储为JSON字符串
    file_drop TEXT,                               -- 档案掉落点，存储为JSON字符串
    activity_drop TEXT,                           -- 活动掉落点，存储为JSON字符串
    cute TEXT,                                    -- 默认Q版立绘
	tech_per_mind REAL,                           -- 每一点心智一提升多少科技点 越高越好
	has_120 BOOLEAN                         	  -- 每一点心智一提升多少科技点 越高越好
);`, TableShip)
	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(fmt.Sprintf("Failed to create table: %v", err))
	}
	fmt.Printf("Table[%s] created successfully\n", LocalDBName)
}

func InsertShipData(db *sql.DB, ships map[string]Ship) error {
	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // 如果出了问题，可以回滚

	// 插入每艘船
	insertCnt := int64(0)
	for _, ship := range ships {
		// 序列化数组字段为 JSON 字符串
		clothingJSON, _ := json.Marshal(ship.Clothing)
		typeJSON, _ := json.Marshal(ship.Type)
		techPointCampJSON, _ := json.Marshal(ship.TechPointCamp)
		techPointJSON, _ := json.Marshal(ship.TechPoint)
		ordinaryDropJSON, _ := json.Marshal(ship.OrdinaryDrop)
		fileDropJSON, _ := json.Marshal(ship.FileDrop)
		activityDropJSON, _ := json.Marshal(ship.ActivityDrop)

		// 构建 INSERT 语句
		query := fmt.Sprintf(`
            INSERT OR IGNORE INTO %s (
                name, avatar, clothing, type, tech_point_count, 
                tech_point_camp, tech_point, rarity, mind_cost, camp, construct_time, 
                install_date, transform_date, ordinary_drop, file_drop, activity_drop, cute, tech_per_mind, has_120
            ) VALUES (
                ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
            )`, TableShip)

		// 执行插入操作
		res, err := tx.Exec(query,
			ship.Name,
			ship.Avatar,
			string(clothingJSON),
			string(typeJSON),
			ship.TechPointCount,
			string(techPointCampJSON),
			string(techPointJSON),
			ship.Rarity,
			ship.MindCost,
			ship.Camp,
			ship.ConstructTime,
			ship.InstallDate,
			ship.TransformDate,
			string(ordinaryDropJSON),
			string(fileDropJSON),
			string(activityDropJSON),
			ship.Cute,
			ship.TechPerMind,
			ship.Has120,
		)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		insertCnt += affected
	}
	log.Info("insert to db count", insertCnt)

	// 提交事务
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func LoadShips(selectSql string) (map[string]Ship, error) {
	retryTimes := 3
retry:
	db, err := database.GetDB(LocalDBName)
	if err != nil {
		err = GetAllShips()
		if err != nil {
			return nil, err
		}
		retryTimes--
		if retryTimes >= 0 {
			goto retry
		}
	}
	defer db.Close()

	// 创建一个空的map来存储所有船
	ships := make(map[string]Ship)

	// 查询所有船只的数据
	rows, err := db.Query(selectSql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var ship Ship
		var id int
		var clothingJSON, typeJSON, techPointCampJSON, techPointJSON, ordinaryDropJSON, fileDropJSON, activityDropJSON string

		// 扫描每一行的字段
		err := rows.Scan(
			&id,
			&ship.Name,
			&ship.Avatar,
			&clothingJSON,
			&typeJSON,
			&ship.TechPointCount,
			&techPointCampJSON,
			&techPointJSON,
			&ship.Rarity,
			&ship.MindCost,
			&ship.Camp,
			&ship.ConstructTime,
			&ship.InstallDate,
			&ship.TransformDate,
			&ordinaryDropJSON,
			&fileDropJSON,
			&activityDropJSON,
			&ship.Cute,
			&ship.TechPerMind,
			&ship.Has120,
		)
		if err != nil {
			return nil, err
		}

		// 反序列化 JSON 字符串为对应的数组或切片
		if err := json.Unmarshal([]byte(clothingJSON), &ship.Clothing); err != nil {
			return nil, fmt.Errorf("failed to unmarshal clothing: %v", err)
		}
		if err := json.Unmarshal([]byte(typeJSON), &ship.Type); err != nil {
			return nil, fmt.Errorf("failed to unmarshal type: %v", err)
		}
		if err := json.Unmarshal([]byte(techPointCampJSON), &ship.TechPointCamp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tech_point_camp: %v", err)
		}
		if err := json.Unmarshal([]byte(techPointJSON), &ship.TechPoint); err != nil {
			return nil, fmt.Errorf("failed to unmarshal tech_point: %v", err)
		}
		if err := json.Unmarshal([]byte(ordinaryDropJSON), &ship.OrdinaryDrop); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ordinary_drop: %v", err)
		}
		if err := json.Unmarshal([]byte(fileDropJSON), &ship.FileDrop); err != nil {
			return nil, fmt.Errorf("failed to unmarshal file_drop: %v", err)
		}
		if err := json.Unmarshal([]byte(activityDropJSON), &ship.ActivityDrop); err != nil {
			return nil, fmt.Errorf("failed to unmarshal activity_drop: %v", err)
		}

		// 将船只信息加入 map，船名作为 key
		ships[ship.Name] = ship
	}

	// 检查遍历过程中是否发生错误
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ships, nil
}

func UpdateShips(updateSql string) error {
	db, err := database.GetDB(LocalDBName)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(updateSql)
	return err
}
