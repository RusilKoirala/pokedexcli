package quest

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type QuestType string

const (
	QuestCatchPokemon  QuestType = "catch_pokemon"
	QuestReachLevel    QuestType = "reach_level"
	QuestDefeatTrainer QuestType = "defeat_trainer"
	QuestVisitLocation QuestType = "visit_location"
)

type QuestStatus string

const (
	QuestActive    QuestStatus = "active"
	QuestCompleted QuestStatus = "completed"
	QuestClaimed   QuestStatus = "claimed"
	QuestLocked    QuestStatus = "locked"
)

type Quest struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Type        QuestType   `json:"type"`
	Target      int         `json:"traget"`
	Progress    int         `json:"progress"`
	Status      QuestStatus `json:"status"`
	RewardXP    int         `json:"reward_xp"`
	RewardLevel int         `json:"reward_level"`
	GiverNPCID  string      `json:"giver_npc_id"`
}

type QuestManager struct {
	Quest map[string]*Quest `json:"quests"`
}

func NewQuestManager() *QuestManager {
	return &QuestManager{
		Quest: make(map[string]*Quest),
	}
}

// adds a new quest
func (qm *QuestManager) AddQuest(quest *Quest) {
	qm.Quest[quest.ID] = quest
}

// update quest progress
func (qm *QuestManager) UpdateProgress(questID string, amount int) bool {
	quest, exists := qm.Quest[questID]
	if !exists || quest.Status != QuestActive {
		return false
	}

	quest.Progress += amount

	if quest.Progress >= quest.Target {
		quest.Status = QuestCompleted
		return true
	}

	return false
}

// all active quests
func (qm *QuestManager) GetActiveQuests() []*Quest {
	active := []*Quest{}
	for _, quest := range qm.Quest {
		if quest.Status == QuestActive {
			active = append(active, quest)
		}
	}
	return active
}

// all completed but unclaimed  quests s
func (qm *QuestManager) GetCompletedQuests() []*Quest {
	completed := []*Quest{}

	for _, quest := range qm.Quest {
		if quest.Status == QuestCompleted {
			completed = append(completed, quest)
		}
	}
	return completed
}

// mark rquest as claimed
func (qm *QuestManager) ClaimReward(questId string) *Quest {
	quest, exits := qm.Quest[questId]
	if !exits || quest.Status != QuestCompleted {
		return nil
	}

	quest.Status = QuestClaimed
	return quest
}

// saves quest
func (qm *QuestManager) Save() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".pokedexcli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(configDir, "quests.json")
	data, err := json.MarshalIndent(qm, "", "  ")
	return os.WriteFile(filePath, data, 0644)
}

// load quest from file
func Load() (*QuestManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return NewQuestManager(), nil
	}

	filePath := filepath.Join(homeDir, ".pokedexcli", "quests.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return InitializeDefaultQuests(), nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return NewQuestManager(), nil
	}

	var qm QuestManager
	if err := json.Unmarshal(data, &qm); err != nil {
		return NewQuestManager(), nil
	}

	return &qm, nil
}

func InitializeDefaultQuests() *QuestManager {
	qm := NewQuestManager()

	qm.AddQuest(&Quest{
		ID:          "oak_first_catch",
		Title:       "Your First Catch",
		Description: "Professor Oak wants you to catch 3 Pokemon",
		Type:        QuestCatchPokemon,
		Target:      3,
		Progress:    0,
		Status:      QuestLocked,
		RewardXP:    100,
		RewardLevel: 0,
		GiverNPCID:  "prof_oak",
	})

	qm.AddQuest(&Quest{
		ID:          "reach_level_5",
		Title:       "Growing Stronger",
		Description: "Reach Level 5 to unlock Virdian Forest",
		Type:        QuestReachLevel,
		Target:      5,
		Progress:    0,
		Status:      QuestLocked,
		RewardXP:    0,
		RewardLevel: 5,
		GiverNPCID:  "prof_oak",
	})

	qm.AddQuest(&Quest{
		ID:          "bug_catcher_challenge",
		Title:       "Bug Catcher's Challenge",
		Description: "Defeat Bug Catcher Jimmy",
		Type:        QuestDefeatTrainer,
		Target:      1,
		Progress:    0,
		Status:      QuestLocked,
		RewardXP:    150,
		RewardLevel: 0,
		GiverNPCID:  "bug_catcher",
	})

	return qm
}

func (qm *QuestManager) OnCatchPokemon() []string {
	completed := []string{}
	for _, quest := range qm.Quest {
		if quest.Type == QuestCatchPokemon && quest.Status == QuestActive {
			if qm.UpdateProgress(quest.ID, 1) {
				completed = append(completed, quest.Title)
			}
		}
	}
	return completed
}

// level quest
func (qm *QuestManager) OnLevelUp(newLevel int) []string {
	completed := []string{}

	for _, quest := range qm.Quest {
		if quest.Type == QuestReachLevel && quest.Status == QuestActive {
			quest.Progress = newLevel
			if quest.Progress >= quest.Target {
				quest.Status = QuestCompleted
				completed = append(completed, quest.Title)
			}
		}
	}
	return completed
}

// trainer defeat quests
func (qm *QuestManager) OnDefeatTrainer(trainerID string) []string {
	completed := []string{}
	for _, quest := range qm.Quest {
		if quest.Type == QuestDefeatTrainer && quest.Status == QuestActive && quest.GiverNPCID == trainerID {
			if qm.UpdateProgress(quest.ID, 1) {
				completed = append(completed, quest.Title)
			}
		}
	}
	return completed
}

// UnlockQuest gives the player ONE quest from this NPC.
// Rules:
//   - If the player already has any active quest = do nothing (busy)
//   - Otherwise unlock only the FIRST locked quest this NPC has
func (qm *QuestManager) UnlockQuest(npcID string) []string {
	// Block if player is already on a quest (from any NPC)
	if len(qm.GetActiveQuests()) > 0 {
		return []string{}
	}

	// Unlock only ONE quest from this NPC
	for _, quest := range qm.Quest {
		if quest.GiverNPCID == npcID && quest.Status == QuestLocked {
			quest.Status = QuestActive
			return []string{quest.Title} // stop after first one
		}
	}
	return []string{}
}
