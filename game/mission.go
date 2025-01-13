package game

import (
	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
)

func NewMission() *sro.MissionBin {
	return &sro.MissionBin{}
}

func GetMissionBin(s *enter.Session) *sro.MissionBin {
	bin := GetPlayerBin(s)
	if bin == nil {
		return nil
	}
	if bin.MissionBin == nil {
		bin.MissionBin = NewMission()
	}
	return bin.MissionBin
}

func GetTutorialList(s *enter.Session) map[int64]bool {
	bin := GetMissionBin(s)
	if bin == nil {
		return nil
	}
	return bin.GetTutorialList()
}

func FinishTutorial(s *enter.Session, tutorialIdS []int64) bool {
	bin := GetMissionBin(s)
	if bin == nil {
		return false
	}
	if bin.TutorialList == nil {
		bin.TutorialList = make(map[int64]bool)
	}
	for _, t := range tutorialIdS {
		bin.TutorialList[t] = true
	}
	return true
}
