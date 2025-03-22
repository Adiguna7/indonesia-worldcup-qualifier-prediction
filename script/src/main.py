from typing import NamedTuple, Tuple
from collections import defaultdict
from enum import Enum, auto


import numpy as np

JAPAN_ELO = 1888
AUSTRALIA_ELO = 1718
INDONESIA_ELO = 1317
SAUDI_ARABIA_ELO = 1535
BAHRAIN_ELO = 1528
CHINA_ELO = 1422

TARGET_TEAM = "idn"


class MatchProbability(NamedTuple):
    home_win: float
    draw: float
    away_win: float


class MatchStatus(Enum):
    HOME_WIN = auto()
    AWAY_WIN = auto()
    DRAW = auto()


class MatchResult(NamedTuple):
    home_team: str
    away_team: str
    status: MatchStatus


MATCH_LEFT = [
    ("aus", "idn"),
    ("jpn", "bhr"),
    ("sau", "chn"),
    ("jpn", "sau"),
    ("chn", "aus"),
    ("idn", "bhr"),
    ("idn", "chn"),
    ("aus", "jpn"),
    ("bhr", "sau"),
    ("jpn", "idn"),
    ("chn", "bhr"),
]

ELO_MAPPING = {
    "jpn": JAPAN_ELO,
    "aus": AUSTRALIA_ELO,
    "idn": INDONESIA_ELO,
    "sau": SAUDI_ARABIA_ELO,
    "bhr": BAHRAIN_ELO,
    "chn": CHINA_ELO
}

FIFA_RANKS = {
    "jpn": 13,
    "aus": 43,
    "idn": 134,
    "sau": 75,
    "bhr": 77,
    "chn": 98,
}

INITIAL_POINTS = {
    "jpn": 16,
    "aus": 7,
    "idn": 6,
    "sau": 6,
    "bhr": 6,
    "chn": 6
}

MATCH_HISTORY: list[MatchResult] = []

# populate the match history with the last 5 matches

# Indonesia
MATCH_HISTORY.extend([
    MatchResult("idn", "sau", MatchStatus.HOME_WIN),
    MatchResult("idn", "jpn", MatchStatus.AWAY_WIN),
    MatchResult("chn", "idn", MatchStatus.HOME_WIN),
    MatchResult("bhr", "idn", MatchStatus.DRAW),
    MatchResult("idn", "aus", MatchStatus.DRAW),
])

# Japan
MATCH_HISTORY.extend([
    MatchResult("chn", "jpn", MatchStatus.AWAY_WIN),
    MatchResult("idn", "jpn", MatchStatus.AWAY_WIN),
    MatchResult("jpn", "aus", MatchStatus.DRAW),
    MatchResult("sau", "jpn", MatchStatus.AWAY_WIN),
    MatchResult("bhr", "jpn", MatchStatus.AWAY_WIN),
])

# Saudi Arabia
MATCH_HISTORY.extend([
    MatchResult("idn", "sau", MatchStatus.HOME_WIN),
    MatchResult("aus", "sau", MatchStatus.DRAW),
    MatchResult("sau", "bhr", MatchStatus.DRAW),
    MatchResult("sau", "jpn", MatchStatus.AWAY_WIN),
    MatchResult("chn", "sau", MatchStatus.AWAY_WIN),
])

# China
MATCH_HISTORY.extend([
    MatchResult("chn", "jpn", MatchStatus.AWAY_WIN),
    MatchResult("bhr", "chn", MatchStatus.AWAY_WIN),
    MatchResult("chn", "idn", MatchStatus.HOME_WIN),
    MatchResult("aus", "chn", MatchStatus.HOME_WIN),
    MatchResult("chn", "sau", MatchStatus.AWAY_WIN)
])

# Bahrain
MATCH_HISTORY.extend([
    MatchResult("bhr", "aus", MatchStatus.DRAW),
    MatchResult("bhr", "chn", MatchStatus.DRAW),
    MatchResult("sau", "bhr", MatchStatus.DRAW),
    MatchResult("bhr", "idn", MatchStatus.DRAW),
    MatchResult("bhr", "jpn", MatchStatus.AWAY_WIN)
])

# Australia
MATCH_HISTORY.extend([
    MatchResult("bhr", "aus", MatchStatus.DRAW),
    MatchResult("aus", "sau", MatchStatus.DRAW),
    MatchResult("jpn", "aus", MatchStatus.DRAW),
    MatchResult("aus", "chn", MatchStatus.HOME_WIN),
    MatchResult("idn", "aus", MatchStatus.DRAW),
])


NUMS_OF_SIMULATION = 10_000


def calculate_match_probability(match: Tuple[str, str]) -> MatchProbability:
    home_team, away_team = match

    home_elo = ELO_MAPPING[home_team]
    away_elo = ELO_MAPPING[away_team]

    p_home_win = calculate_elo_probability(home_elo, away_elo)
    p_draw = 0.25
    p_away_win = 1 - (p_home_win + p_draw)

    home_boost, away_boost = calculate_rank_probability(match)

    p_home_win += home_boost
    p_away_win += away_boost

    home_boost, away_boost = calculate_recent_match_boost(match)
    p_home_win += home_boost
    p_away_win += away_boost

    # home advantages
    p_home_win += 0.10

    # normalization
    total = p_home_win + p_draw + p_away_win
    p_home_win /= total
    p_draw /= total
    p_away_win /= total

    return MatchProbability(p_home_win, p_draw, p_away_win)


def calculate_elo_probability(home_elo: int, away_elo: int) -> float:
    return 1 / (1 + pow(10, (away_elo - home_elo) / 400))


def calculate_boost_probability(
    match: Tuple[str, str],
    base_boost: float,
    difference: float,
    max_difference: float
) -> Tuple[float, float]:
    if difference < 0:
        boost = 0.05 * (difference / max_difference)
        return (boost, 0.0)

    boost = 0.05 * (-difference / max_difference)
    return (0.0, boost)


def calculate_rank_probability(match: Tuple[str, str]):
    home_team, away_team = match

    # assume the full boost of higher fifa rating is 5%
    # but this adjustable by how far both of team rank
    max_fifa_rank_difference = \
        max(FIFA_RANKS.values()) - min(FIFA_RANKS.values())

    fifa_rank_difference = FIFA_RANKS[home_team] - FIFA_RANKS[away_team]
    return calculate_boost_probability(
        match=match,
        base_boost=0.05,
        difference=fifa_rank_difference,
        max_difference=max_fifa_rank_difference
    )


def calculate_recent_match_boost(
    match: Tuple[str, str]
) -> Tuple[float, float]:
    home_team, away_team = match
    history: dict[str, list[MatchResult]] = {}

    for past_match in MATCH_HISTORY:
        if past_match.home_team not in history:
            history[past_match.home_team] = [past_match]
        else:
            history[past_match.home_team].append(past_match)

        if past_match.away_team not in history:
            history[past_match.away_team] = [past_match]
        else:
            history[past_match.away_team].append(past_match)

    home_boost = 0.0
    away_boost = 0.0

    max_elo_diff = max(ELO_MAPPING.values()) - min(ELO_MAPPING.values())

    def calculate_boost_penalty(
        team: str, opponent: str, status: MatchStatus
    ) -> float:
        elo_diff = abs(ELO_MAPPING[team] - ELO_MAPPING[opponent])
        weight = elo_diff / max_elo_diff
        base_boost = 0.05 * weight

        if ELO_MAPPING[team] < ELO_MAPPING[opponent]:
            if status in {MatchStatus.HOME_WIN, MatchStatus.AWAY_WIN}:
                return base_boost
            elif status == MatchStatus.DRAW:
                return base_boost * 0.6
        else:
            if status in {MatchStatus.HOME_WIN, MatchStatus.AWAY_WIN}:
                return 0.0
            elif status == MatchStatus.DRAW:  # draw vs weaker team
                return -base_boost * 0.5  # 50% of max penalty
            elif (
                    team == past_match.home_team and
                    status == MatchStatus.AWAY_WIN
                ) or \
                (
                    team == past_match.away_team and
                    status == MatchStatus.HOME_WIN
                    ):
                return -base_boost

        return 0.0

    for past_match in history[home_team]:
        opponent = past_match.away_team if past_match.home_team == \
            home_team else past_match.home_team
        home_boost += calculate_boost_penalty(
            home_team,
            opponent,
            past_match.status
        )

    for past_match in history[away_team]:
        opponent = past_match.away_team if past_match.home_team == \
            away_team else past_match.home_team
        away_boost += calculate_boost_penalty(
            away_team,
            opponent,
            past_match.status
        )

    return home_boost, away_boost


if __name__ == "__main__":
    match_probability: dict[Tuple[str, str], MatchProbability] = {}
    direct_qualified = 0
    pass_to_next_round = 0

    for match in MATCH_LEFT:
        match_probability[match] = \
            calculate_match_probability(match)

    for i in range(NUMS_OF_SIMULATION):
        current_point = defaultdict(int, INITIAL_POINTS)

        for match in MATCH_LEFT:
            home_team, away_team = match
            match_result = np.random.rand()

            probability = match_probability[match]

            if match_result < probability.home_win:
                current_point[home_team] += 3
            elif match_result < probability.home_win + probability.draw:
                current_point[home_team] += 1
                current_point[away_team] += 1
            else:
                current_point[away_team] += 3

        sorted_point = sorted(
            current_point.items(), key=lambda x: x[1], reverse=True
        )

        for index, (team, _) in enumerate(sorted_point):
            if team == TARGET_TEAM:
                if index < 2:
                    direct_qualified += 1
                elif index < 4:
                    pass_to_next_round += 1
                break

        print(f"================ simulation {i + 1} ================")
        print(sorted_point)
        print("=====================================================\n\n")

    print(f"chance of {TARGET_TEAM} directly qualified for world cup {direct_qualified / NUMS_OF_SIMULATION:.2%}")  # noqa
    print(f"chance of {TARGET_TEAM} pass to the next round {pass_to_next_round / NUMS_OF_SIMULATION:.2%}")  # noqa
