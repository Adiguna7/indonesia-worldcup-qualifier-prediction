package main

import (
	"fmt"
	"log"
	"maps"
	"math"
	"math/rand"
	"sort"
)

// TODO: UPDATE ELO FROM ELO RATING WEB
const JAPAN_ELO = 1888
const AUSTRALIA_ELO = 1718
const INDONESIA_ELO = 1317
const SAUDI_ARABIA_ELO = 1535
const BAHRAIN_ELO = 1528
const CHINA_ELO = 1422

const TARGET_TEAM = "idn"
const NUMS_OF_SIMULATION = 10000

type MatchProbability struct {
	homeWin float64
	draw    float64
	awayWin float64
}

type MatchStatus int

const (
	HOME_WIN MatchStatus = iota
	DRAW
	AWAY_WIN
	NOT_PLAYED_YET
)

type TeamMatchResult int

const (
	WIN TeamMatchResult = iota
	DRAW_RESULT
	LOSE
)

type Match struct {
	homeTeam string
	awayTeam string
	status   MatchStatus
}

type ProbabilityBoost struct {
	homeTeam float64
	awayTeam float64
}

var matchesLeft []Match = []Match{
	{homeTeam: "aus", awayTeam: "idn", status: NOT_PLAYED_YET},
	{homeTeam: "jpn", awayTeam: "bhr", status: NOT_PLAYED_YET},
	{homeTeam: "sau", awayTeam: "chn", status: NOT_PLAYED_YET},
	{homeTeam: "jpn", awayTeam: "sau", status: NOT_PLAYED_YET},
	{homeTeam: "chn", awayTeam: "aus", status: NOT_PLAYED_YET},
	{homeTeam: "idn", awayTeam: "bhr", status: NOT_PLAYED_YET},
	{homeTeam: "idn", awayTeam: "chn", status: NOT_PLAYED_YET},
	{homeTeam: "aus", awayTeam: "jpn", status: NOT_PLAYED_YET},
	{homeTeam: "bhr", awayTeam: "sau", status: NOT_PLAYED_YET},
	{homeTeam: "jpn", awayTeam: "idn", status: NOT_PLAYED_YET},
	{homeTeam: "chn", awayTeam: "bhr", status: NOT_PLAYED_YET},
}

var eloMapping map[string]int = map[string]int{
	"jpn": JAPAN_ELO,
	"aus": AUSTRALIA_ELO,
	"idn": INDONESIA_ELO,
	"sau": SAUDI_ARABIA_ELO,
	"bhr": BAHRAIN_ELO,
	"chn": CHINA_ELO,
}

var fifaRanks map[string]int = map[string]int{
	"jpn": 13,
	"aus": 43,
	"idn": 134,
	"sau": 75,
	"bhr": 77,
	"chn": 98,
}

var maxRank, minRank = findMinMax(fifaRanks)

var initialPoints map[string]int = map[string]int{
	"jpn": 16,
	"aus": 7,
	"idn": 6,
	"sau": 6,
	"bhr": 6,
	"chn": 6,
}

var matchHistory = [...]Match{
	{homeTeam: "idn", awayTeam: "sau", status: HOME_WIN},
	{homeTeam: "idn", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "chn", awayTeam: "idn", status: HOME_WIN},
	{homeTeam: "bhr", awayTeam: "idn", status: DRAW},
	{homeTeam: "idn", awayTeam: "aus", status: DRAW},
	{homeTeam: "chn", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "idn", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "jpn", awayTeam: "aus", status: DRAW},
	{homeTeam: "sau", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "bhr", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "idn", awayTeam: "sau", status: HOME_WIN},
	{homeTeam: "aus", awayTeam: "sau", status: DRAW},
	{homeTeam: "sau", awayTeam: "bhr", status: DRAW},
	{homeTeam: "sau", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "chn", awayTeam: "sau", status: AWAY_WIN},
	{homeTeam: "chn", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "bhr", awayTeam: "chn", status: AWAY_WIN},
	{homeTeam: "chn", awayTeam: "idn", status: HOME_WIN},
	{homeTeam: "aus", awayTeam: "chn", status: HOME_WIN},
	{homeTeam: "chn", awayTeam: "sau", status: AWAY_WIN},
	{homeTeam: "bhr", awayTeam: "aus", status: DRAW},
	{homeTeam: "bhr", awayTeam: "chn", status: DRAW},
	{homeTeam: "sau", awayTeam: "bhr", status: DRAW},
	{homeTeam: "bhr", awayTeam: "idn", status: DRAW},
	{homeTeam: "bhr", awayTeam: "jpn", status: AWAY_WIN},
	{homeTeam: "bhr", awayTeam: "aus", status: DRAW},
	{homeTeam: "aus", awayTeam: "sau", status: DRAW},
	{homeTeam: "jpn", awayTeam: "aus", status: DRAW},
	{homeTeam: "aus", awayTeam: "chn", status: HOME_WIN},
	{homeTeam: "idn", awayTeam: "aus", status: DRAW},
}

type Team struct {
	name  string
	point int
}

func main() {
	matchProbability := map[Match]MatchProbability{}
	directQualified := 0
	passToNextRound := 0

	for _, match := range matchesLeft {
		probabilityResult, err := calculateMatchProbability(match)

		if err != nil {
			log.Fatalf("error: %v", err)
		}

		matchProbability[match] = probabilityResult
	}

	for range NUMS_OF_SIMULATION {
		currentPoint := make(map[string]int)
		maps.Copy(currentPoint, initialPoints)

		for _, match := range matchesLeft {
			homeTeam, awayTeam := match.homeTeam, match.awayTeam
			matchResult := rand.Float64()

			probability := matchProbability[match]

			if matchResult < probability.homeWin {
				currentPoint[homeTeam] += 3
			} else if matchResult < probability.homeWin+probability.draw {
				currentPoint[homeTeam] += 1
				currentPoint[awayTeam] += 1
			} else {
				currentPoint[awayTeam] += 3
			}
		}

		var standings []Team
		for key, value := range currentPoint {
			standings = append(standings, Team{
				name:  key,
				point: value,
			})
		}

		sort.Slice(standings, func(i, j int) bool {
			return standings[i].point > standings[j].point
		})

		for index, team := range standings {
			if team.name == TARGET_TEAM {
				if index < 2 {
					directQualified += 1
				} else if index < 4 {
					passToNextRound += 1
				}
				break
			}
		}
	}
	fmt.Printf("Chance of %s directly qualifying for World Cup: %.2f%%\n", TARGET_TEAM, float64(directQualified)/float64(NUMS_OF_SIMULATION)*100)
	fmt.Printf("Chance of %s passing to the next round: %.2f%%\n", TARGET_TEAM, float64(passToNextRound)/float64(NUMS_OF_SIMULATION)*100)
}

func calculateMatchProbability(match Match) (MatchProbability, error) {
	clampFunction := func(value float64) float64 {
		if value < 0 {
			return 0
		}

		return value
	}

	homeTeam := match.homeTeam
	awayTeam := match.awayTeam

	homeElo := eloMapping[homeTeam]
	awayElo := eloMapping[awayTeam]

	pHomeWin := calculateEloProbability(homeElo, awayElo)
	pDraw := 0.25
	pAwayWin := 1 - (pHomeWin + pDraw)

	boost := calculateRankProbability(match)

	pHomeWin += boost.homeTeam
	pAwayWin += boost.awayTeam

	probabilityBoost := calculateRecentMatchBoost(match)

	pHomeWin += probabilityBoost.homeTeam
	pAwayWin += probabilityBoost.awayTeam

	// home advantages
	pHomeWin += 0.10

	// clamp value
	pHomeWin = clampFunction(pHomeWin)
	pAwayWin = clampFunction(pAwayWin)
	pDraw = clampFunction(pDraw)

	//  normalization
	total := pHomeWin + pDraw + pAwayWin
	pHomeWin /= total
	pDraw /= total
	pAwayWin /= total

	return MatchProbability{
		homeWin: pHomeWin,
		draw:    pDraw,
		awayWin: pAwayWin,
	}, nil
}

func calculateEloProbability(homeElo int, awayElo int) float64 {
	return 1 / (1 + math.Pow(float64(10), float64(awayElo-homeElo)/400))
}

func calculateRankProbability(match Match) ProbabilityBoost {
	homeTeam := match.homeTeam
	awayTeam := match.awayTeam

	// assume the full boost of higher fifa rating is 5%
	// but this adjustable by how far both of team rank
	maxRankDifference := maxRank - minRank

	rankDifference := fifaRanks[homeTeam] - fifaRanks[awayTeam]
	if rankDifference < 0 {
		boost := 0.05 * (float64(-rankDifference) / float64(maxRankDifference))
		return ProbabilityBoost{
			homeTeam: 0.0,
			awayTeam: boost,
		}
	}

	boost := 0.05 * (float64(rankDifference) / float64(maxRankDifference))
	return ProbabilityBoost{
		homeTeam: boost,
		awayTeam: 0.0,
	}
}

func calculateRecentMatchBoost(match Match) ProbabilityBoost {
	homeTeam, awayTeam := match.homeTeam, match.awayTeam
	history := map[string][]Match{}

	for _, pastMatch := range matchHistory {
		pastMatchHome := pastMatch.homeTeam
		pastMatchAway := pastMatch.awayTeam

		if _, exist := history[pastMatchHome]; !exist {
			history[pastMatchHome] = []Match{
				pastMatch,
			}
		} else {
			history[pastMatchHome] = append(history[pastMatchHome], pastMatch)
		}

		if _, exist := history[pastMatchAway]; !exist {
			history[pastMatchAway] = []Match{
				pastMatch,
			}
		} else {
			history[pastMatchAway] = append(history[pastMatchAway], pastMatch)
		}
	}

	homeBoost := 0.0
	awayBoost := 0.0

	maxEloDiff, _ := findMinMax(eloMapping)

	calculateBoostPenalty := func(
		team string,
		opponent string,
		status TeamMatchResult,
	) float64 {

		eloDiff := math.Abs(float64(eloMapping[team] - eloMapping[opponent]))
		weight := eloDiff / float64(maxEloDiff)
		baseBoost := 0.05 * weight

		if eloMapping[team] < eloMapping[opponent] {
			if status == WIN {
				return baseBoost
			} else if status == DRAW_RESULT {
				return baseBoost * 0.6
			}
		} else {
			if status == WIN {
				return 0.0
			} else if status == DRAW_RESULT { // draw vs weaker team
				return -baseBoost * 0.5
			} else if status == LOSE {
				return -baseBoost
			}
		}

		return 0.0
	}
	for _, pastMatch := range history[homeTeam] {
		var team string
		var opponent string
		var status TeamMatchResult

		if pastMatch.homeTeam == homeTeam {
			team = pastMatch.homeTeam
			opponent = pastMatch.awayTeam

			if pastMatch.status == HOME_WIN {
				status = WIN
			} else if pastMatch.status == DRAW {
				status = DRAW_RESULT
			} else {
				status = LOSE
			}

		} else {
			team = pastMatch.awayTeam
			opponent = pastMatch.homeTeam

			if pastMatch.status == AWAY_WIN {
				status = WIN
			} else if pastMatch.status == DRAW {
				status = DRAW_RESULT
			} else {
				status = LOSE
			}
		}

		homeBoost += calculateBoostPenalty(team, opponent, status)
	}

	for _, pastMatch := range history[awayTeam] {
		var team string
		var opponent string
		var status TeamMatchResult

		if pastMatch.awayTeam == awayTeam {
			team = pastMatch.awayTeam
			opponent = pastMatch.homeTeam

			if pastMatch.status == AWAY_WIN {
				status = WIN
			} else if pastMatch.status == DRAW {
				status = DRAW_RESULT
			} else {
				status = LOSE
			}

		} else {
			team = pastMatch.homeTeam
			opponent = pastMatch.awayTeam

			if pastMatch.status == AWAY_WIN {
				status = WIN
			} else if pastMatch.status == DRAW {
				status = DRAW_RESULT
			} else {
				status = LOSE
			}
		}

		awayBoost += calculateBoostPenalty(
			team,
			opponent,
			status,
		)
	}
	return ProbabilityBoost{
		homeTeam: homeBoost,
		awayTeam: awayBoost,
	}
}

func findMinMax(scores map[string]int) (maxValue int, minValue int) {
	maxValue = math.MinInt
	minValue = math.MaxInt

	for _, score := range scores {
		if score > maxValue {
			maxValue = score
		}
		if score < minValue {
			minValue = score
		}
	}
	return
}
