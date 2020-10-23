package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/DSU-DefSec/mew/checks"
)

func validateString(input string) bool {
	if input == "" {
		return false
	}
	validationString := `^[a-zA-Z0-9-_]+$`
	inputValidation := regexp.MustCompile(validationString)
	return inputValidation.MatchString(input)
}

func (m *config) validateTeam(teamName string) (teamData, error) {
	if !validateString(teamName) {
		return teamData{}, errors.New("team string contains illegal characters")
	}
	team, err := m.getTeam(teamName)
	if err != nil {
		return team, err
	}
	return team, nil
}

func (m *config) getCheck(checkName string) (checks.Check, error) {
	for _, check := range m.CheckList {
		fmt.Println("check.Name is", check.FetchName(), "checkNAme is", checkName)
		if check.FetchName() == checkName {
			return check, nil
		}
	}
	return checks.Web{}, errors.New("check not found")
}

func (m *config) getTeam(teamName string) (teamData, error) {
	for _, team := range m.Team {
		if team.Name == teamName {
			return team, nil
		}
	}
	return teamData{}, errors.New("team not found")
}

func (m *config) validateTeamIndex(teamName string, teamIndex string) (teamData, error) {
	if len(teamIndex) < 5 {
		return teamData{}, errors.New("team name had invalid length")
	}
	teamNum := teamIndex[4:]
	index, err := strconv.Atoi(teamNum)
	if err != nil {
		return teamData{}, errors.New("team index string is not a number")
	}
	team, err := m.getTeamByIndex(index)
	if err != nil {
		return team, err
	}
	if team.Name != teamName {
		return team, errors.New("unauthorized team")
	}
	return team, nil
}

func (m *config) getTeamByIndex(index int) (teamData, error) {
	if index >= 0 && index < len(m.Team) {
		return m.Team[index], nil
	}
	return teamData{}, errors.New("team index not found")
}