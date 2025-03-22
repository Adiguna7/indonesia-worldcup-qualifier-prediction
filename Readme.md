# Indonesia World Cup Qualifier Prediction Using Monte Carlo and Elo-Based Rating

## Overview
This script simulates the upcoming Asian World Cup qualifiers for Group C using Elo ratings, FIFA rankings, and recent match history to estimate the probabilities of match outcomes. It runs Monte Carlo simulations to determine the likelihood of different teams qualifying based on the given match schedule.

## Assumptions
- The probability of a match outcome is based on:
  - **Elo rating**
  - **FIFA ranking**
  - **Recent match history**:
    - If a lower-rated team beats a higher-rated team in the last five matches, they receive a **5% boost**.
    - If a lower-rated team draws against a higher-rated team, they receive a **5% boost**.
    - All boosts are adjusted using a **scaling factor** based on the Elo rating difference.
  - **Home advantage**: The home team receives a **10% boost**.
- Each minute of the match is simulated independently, without considering momentum shifts.
- The model does not account for injuries, red cards, weather conditions, or referee decisions.
- Team strengths remain constant throughout the match and do not adapt dynamically.
- Each minute is treated equally, without considering fatigue effects.

## Features
- Uses Elo ratings to calculate the probability of a home win, draw, or away win.
- Adjusts probabilities based on FIFA rankings.
- Incorporates recent match history to refine performance-based adjustments.
- Simulates **10,000 tournament scenarios** to evaluate possible outcomes.
- Computes the likelihood of **direct qualification** or **advancing to the next round**.

## Dependencies
- Python 3.10+
- NumPy

## How It Works
1. **Match Setup**:
   - Defines Elo ratings for six teams in Group C (Japan, Australia, Indonesia, Saudi Arabia, Bahrain, China).
   - Establishes initial points for each team (based on current standings).
   - Lists upcoming matches.

2. **Probability Calculation**:
   - Determines win probabilities using the Elo rating system.
   - Adjusts for FIFA ranking differences.
   - Applies recent match history to fine-tune probabilities.

3. **Simulation**:
   - Runs **10,000 simulations** of the remaining matches.
   - Updates the standings based on simulated match results.
   - Tracks the number of times **Indonesia qualifies directly** or **advances to the next round**.

## Running the Simulation
To execute the simulation, run:
only first time
```sh
poetry install
```
```sh
python src/main.py
```

## Improvement
- Implement a proper scaling factor, such as a sigmoid function.
- Consider goal counts from previous matches as part of the performance boost calculation.
- Implement web base to animate the simulation
- Scrape match result from web instead hardcoded

## Author
Suryo Adiguna

