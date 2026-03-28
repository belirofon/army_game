#!/usr/bin/env python3
"""
Balance Simulator for Army Game (Армейка)

Monte Carlo simulation to test game balance before implementation.
Tests: death rate, average turn reached, MOR changes, etc.

Usage:
    python simulator.py --runs 1000
    python simulator.py --runs 5000 --seed 42 --output csv
    python simulator.py --verbose --runs 100

Requirements:
    Python 3.10+
    stdlib only
"""

import argparse
import json
import random
import statistics
import sys
from dataclasses import dataclass, field
from datetime import datetime
from typing import Optional, TypedDict


class EventChoice(TypedDict, total=False):
    id: str
    text: str
    probability: float
    success_effects: dict[str, int]
    partial_effects: dict[str, int]
    failure_effects: dict[str, int]


class EventData(TypedDict, total=False):
    id: str
    type: str
    difficulty: int
    narrative: str
    choices: list[EventChoice]


@dataclass
class PlayerState:
    str_stat: int = 50
    end: int = 50
    agi: int = 50
    mor: int = 50
    disc: int = 0
    turn: int = 1
    flags: list[str] = field(default_factory=list)
    history: list[dict] = field(default_factory=list)

    def clone(self) -> 'PlayerState':
        return PlayerState(
            str_stat=self.str_stat, end=self.end, agi=self.agi,
            mor=self.mor, disc=self.disc, turn=self.turn,
            flags=self.flags.copy(), history=self.history.copy()
        )

    def apply_effect(self, stat: str, delta: int) -> None:
        if stat == 'str':
            self.str_stat = max(0, min(100, self.str_stat + delta))
        elif stat == 'end':
            self.end = max(0, min(100, self.end + delta))
        elif stat == 'agi':
            self.agi = max(0, min(100, self.agi + delta))
        elif stat == 'mor':
            self.mor = max(0, min(100, self.mor + delta))
        elif stat == 'disc':
            self.disc = max(-100, min(100, self.disc + delta))


@dataclass
class Event:
    id: str
    type: str
    difficulty: int
    choices: list[EventChoice]
    narrative: str = ""

    def get_effects(self, outcome: str, choice: EventChoice) -> dict[str, int]:
        if outcome == 'success':
            return choice.get('success_effects', {})
        elif outcome == 'partial':
            return choice.get('partial_effects', {})
        else:
            return choice.get('failure_effects', {})


class BalanceSimulator:
    """Core simulation engine."""

    EVENT_TYPES = ['ROUTINE', 'SOCIAL', 'INSPECTION', 'INFORMAL', 'EMERGENCY', 'SAFE']
    STATS = ['str', 'end', 'agi', 'mor', 'disc']

    def __init__(self, events: list[EventData]):
        self.events = [Event(**e) for e in events]
        self.recent_negatives = 0

    def simulate(self) -> dict:
        """Run one complete simulation."""
        player = PlayerState()
        self.recent_negatives = 0

        while True:
            if player.mor <= 0:
                return {
                    'status': 'death',
                    'turn': player.turn,
                    'final_mor': player.mor,
                    'final_disc': player.disc
                }

            if player.turn > 30:
                perfect = player.mor > 50 and -20 <= player.disc <= 20
                return {
                    'status': 'victory',
                    'turn': player.turn,
                    'final_mor': player.mor,
                    'final_disc': player.disc,
                    'perfect': perfect
                }

            difficulty_mod = self._get_difficulty_mod(player.turn)

            if self.recent_negatives >= 2:
                event = self._select_safe_event()
            else:
                event = self._select_event(player, difficulty_mod)

            choice = random.choice(event.choices)
            outcome = self._roll_outcome(choice, event.difficulty, difficulty_mod)
            effects = event.get_effects(outcome, choice)

            for stat, delta in effects.items():
                player.apply_effect(stat, delta)

            if 'mor' in effects and effects['mor'] < 0:
                self.recent_negatives += 1
            elif 'mor' in effects and effects['mor'] > 0:
                self.recent_negatives = 0

            player.turn += 1
            player.history.append({
                'event_id': event.id,
                'event_type': event.type,
                'choice_id': choice['id'],
                'outcome': outcome,
                'effects': effects
            })

    def _get_difficulty_mod(self, turn: int) -> float:
        """Return difficulty modifier based on turn number.
        
        Higher value = lower probability = harder
        """
        if turn <= 10:
            return 0.60
        elif turn <= 20:
            return 0.85
        else:
            return 1.05

    def _select_event(self, player: PlayerState, difficulty_mod: float) -> Event:
        """Select event using weighted random selection."""
        weights = []
        for e in self.events:
            w = 1.0

            if e.type == 'SAFE' and player.mor < 30:
                w *= 1.5
            elif e.type == 'INSPECTION' and player.disc > 50:
                w *= 1.3
            elif e.type == 'INFORMAL' and player.disc < -50:
                w *= 1.3
            elif e.type == 'ROUTINE' and player.turn <= 10:
                w *= 1.2
            elif e.type == 'EMERGENCY' and player.turn > 20:
                w *= 1.2

            w *= (1.0 + (5 - e.difficulty) * 0.1)

            weights.append(w)

        total = sum(weights)
        probs = [w / total for w in weights]
        return random.choices(self.events, weights=probs)[0]

    def _select_safe_event(self) -> Event:
        """Select a safe event for recovery."""
        safe = [e for e in self.events if e.type == 'SAFE']
        if safe:
            return random.choice(safe)
        return random.choice(self.events)

    def _roll_outcome(self, choice: EventChoice, difficulty: int, difficulty_mod: float) -> str:
        """Roll for outcome based on probability and difficulty.
        
        Lower difficulty_mod = easier (early game)
        Higher difficulty_mod = harder (late game)
        """
        base_prob = choice['probability']

        adjusted_prob = base_prob / difficulty_mod

        difficulty_penalty = (difficulty - 1) * 0.06
        adjusted_prob -= difficulty_penalty
        adjusted_prob = max(0.20, min(0.85, adjusted_prob))

        partial_bonus = 0.35
        roll = random.random()
        if roll < adjusted_prob:
            return 'success'
        elif roll < adjusted_prob + (1 - adjusted_prob) * partial_bonus:
            return 'partial'
        else:
            return 'failure'


class StatisticsAnalyzer:
    """Analyze simulation results."""

    TARGETS = {
        'death_rate': (0.20, 0.35),
        'victory_rate': (0.65, 0.80),
        'perfect_rate': (0.00, 0.05),
        'avg_turn_death': (12, 28),
        'avg_mor_change': (-2.0, -1.2),
    }

    def __init__(self, results: list[dict]):
        self.results = results
        self.deaths = [r for r in results if r['status'] == 'death']
        self.victories = [r for r in results if r['status'] == 'victory']
        self.perfect = [r for r in self.victories if r.get('perfect', False)]

    def get_stats(self) -> dict:
        """Calculate all statistics."""
        total = len(self.results)
        if total == 0:
            return {}

        death_rate = len(self.deaths) / total
        victory_rate = len(self.victories) / total
        perfect_rate = len(self.perfect) / total if self.victories else 0

        death_turns = [r['turn'] for r in self.deaths] if self.deaths else [0]
        avg_turn_death = statistics.mean(death_turns)

        all_mor_changes = []
        for r in self.results:
            mor_change = (r['final_mor'] - 50) / max(1, r['turn'] - 1) if r['turn'] > 1 else 0
            all_mor_changes.append(mor_change)
        avg_mor_change = statistics.mean(all_mor_changes)

        return {
            'total_runs': total,
            'death_count': len(self.deaths),
            'victory_count': len(self.victories),
            'perfect_count': len(self.perfect),
            'death_rate': death_rate,
            'victory_rate': victory_rate,
            'perfect_rate': perfect_rate,
            'avg_turn_death': avg_turn_death,
            'avg_mor_change': avg_mor_change,
            'death_turns': death_turns,
        }

    def check_targets(self, stats: dict) -> dict[str, bool]:
        """Check if metrics are within targets."""
        checks = {}

        checks['death_rate'] = (
            self.TARGETS['death_rate'][0] <= stats['death_rate'] <= self.TARGETS['death_rate'][1]
        )
        checks['victory_rate'] = (
            self.TARGETS['victory_rate'][0] <= stats['victory_rate'] <= self.TARGETS['victory_rate'][1]
        )
        checks['perfect_rate'] = (
            self.TARGETS['perfect_rate'][0] <= stats['perfect_rate'] <= self.TARGETS['perfect_rate'][1]
        )

        return checks


class ReportFormatter:
    """Format output for terminal."""

    @staticmethod
    def format_terminal(stats: dict, checks: dict) -> str:
        """Format statistics as ASCII table."""
        death_rate = stats['death_rate'] * 100
        victory_rate = stats['victory_rate'] * 100
        perfect_rate = stats['perfect_rate'] * 100
        avg_turn = stats.get('avg_turn_death', 0)
        avg_mor = stats.get('avg_mor_change', 0)

        bar_death = ReportFormatter._make_bar(death_rate, 100, 30)
        bar_victory = ReportFormatter._make_bar(victory_rate, 100, 80)
        bar_perfect = ReportFormatter._make_bar(perfect_rate, 20, 8)

        death_ok = checks.get('death_rate', False)
        victory_ok = checks.get('victory_rate', False)
        perfect_ok = checks.get('perfect_rate', False)

        ok_mark = '✓'
        fail_mark = '✗'

        lines = [
            "╔══════════════════════════════════════════════════════════════════════╗",
            "║                    BALANCE SIMULATOR REPORT                          ║",
            "╠══════════════════════════════════════════════════════════════════════╣",
           f"║  Runs: {stats['total_runs']:<6}    Deaths: {stats['death_count']:<4} ╣  Victories: {stats['victory_count']:<4}                      ║",
            "╠══════════════════════════════════════════════════════════════════════╣",
            "║  KEY METRICS                                                         ║",
            "║  ────────────────────────────────────────────────────────────────────║",
           f"║  Death Rate:        {death_rate:5.1f}%  {bar_death}  target: 20-30%  {ok_mark if death_ok else fail_mark}   ║",
           f"║  Victory Rate:      {victory_rate:5.1f}%  {bar_victory}  target: 70-80%  {ok_mark if victory_ok else fail_mark}   ║",
           f"║  Perfect Runs:       {perfect_rate:5.1f}%  {bar_perfect}  target: 5-10%   {ok_mark if perfect_ok else fail_mark}   ║",
            "║  ──────────────────────────────────────────────────────────────────── ║",
           f"║  Avg Death Turn:    {avg_turn:5.1f}                                               ║",
           f"║  Avg MOR/turn:       {avg_mor:+5.2f}     target: -2.5 to -1.5              ║",
            "╠══════════════════════════════════════════════════════════════════════╣",
        ]

        all_ok = all(checks.values())
        status = "BALANCED ✓" if all_ok else "NEEDS TUNING ✗"
        lines.append(f"║  FINAL STATUS: {status:<57}║")
        lines.append("╚══════════════════════════════════════════════════════════════════════╝")

        return '\n'.join(lines)

    @staticmethod
    def _make_bar(value: float, max_val: float, target: float) -> str:
        """Create ASCII progress bar."""
        width = 20
        filled = int((value / max_val) * width)
        target_pos = int((target / max_val) * width)
        target_pos = min(target_pos, width - 1)

        bar = ['░'] * width
        for i in range(filled):
            bar[i] = '█'

        return ''.join(bar)

    @staticmethod
    def format_csv_header() -> str:
        return "total_runs,death_count,victory_count,perfect_count,death_rate,victory_rate,perfect_rate,avg_turn_death,avg_mor_change"

    @staticmethod
    def format_csv_row(stats: dict) -> str:
        return f"{stats['total_runs']},{stats['death_count']},{stats['victory_count']},{stats['perfect_count']},{stats['death_rate']:.4f},{stats['victory_rate']:.4f},{stats['perfect_rate']:.4f},{stats.get('avg_turn_death', 0):.2f},{stats.get('avg_mor_change', 0):.4f}"

    @staticmethod
    def format_json(stats: dict, checks: dict) -> str:
        return json.dumps({
            'statistics': stats,
            'target_checks': checks,
            'all_targets_met': all(checks.values())
        }, indent=2)


def main():
    parser = argparse.ArgumentParser(
        description='Balance Simulator for Army Game',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python simulator.py --runs 1000
  python simulator.py --runs 5000 --seed 42
  python simulator.py --runs 100 --verbose
  python simulator.py --runs 1000 --output csv > results.csv
  python simulator.py --runs 1000 --save results/run1.csv
        """
    )
    parser.add_argument(
        '--runs', '-n', type=int, default=1000,
        help='Number of simulations to run (default: 1000)'
    )
    parser.add_argument(
        '--seed', '-s', type=int, default=None,
        help='Random seed for reproducibility'
    )
    parser.add_argument(
        '--events', '-e', default='events.json',
        help='Path to events JSON file (default: events.json)'
    )
    parser.add_argument(
        '--output', '-o', choices=['terminal', 'csv', 'json'], default='terminal',
        help='Output format (default: terminal)'
    )
    parser.add_argument(
        '--save', type=str, default=None,
        help='Save results to CSV file'
    )
    parser.add_argument(
        '--verbose', '-v', action='store_true',
        help='Print each playthrough result'
    )

    args = parser.parse_args()

    try:
        with open(args.events) as f:
            events_data = json.load(f)
    except FileNotFoundError:
        print(f"Error: Events file '{args.events}' not found", file=sys.stderr)
        sys.exit(1)
    except json.JSONDecodeError as e:
        print(f"Error: Invalid JSON in events file: {e}", file=sys.stderr)
        sys.exit(1)

    if args.seed:
        random.seed(args.seed)

    sim = BalanceSimulator(events_data)
    results = []

    for i in range(args.runs):
        result = sim.simulate()
        results.append(result)

        if args.verbose:
            status = result['status']
            turn = result['turn']
            if status == 'death':
                print(f"  Run {i+1}: DEATH @ turn {turn}")
            elif result.get('perfect'):
                print(f"  Run {i+1}: VICTORY (PERFECT) @ turn {turn}")
            else:
                print(f"  Run {i+1}: VICTORY @ turn {turn}")

    analyzer = StatisticsAnalyzer(results)
    stats = analyzer.get_stats()
    checks = analyzer.check_targets(stats)

    if args.save:
        import os
        import shutil
        os.makedirs(os.path.dirname(args.save) or '.', exist_ok=True)
        
        # Save CSV results
        with open(args.save, 'w') as f:
            f.write(ReportFormatter.format_csv_header() + '\n')
            f.write(ReportFormatter.format_csv_row(stats) + '\n')
        
        # Save metadata
        base_name = args.save.rsplit('.', 1)[0]
        meta_file = base_name + '_meta.json'
        
        metadata = {
            'seed': args.seed,
            'runs': args.runs,
            'events_file': args.events,
            'difficulty_modifier': {
                'turn_1_10': 0.60,
                'turn_11_20': 0.85,
                'turn_21_30': 1.05
            },
            'recovery_trigger': 2,
            'targets': StatisticsAnalyzer.TARGETS,
            'simulation_date': str(datetime.now())
        }
        
        with open(meta_file, 'w') as f:
            json.dump(metadata, f, indent=2)
        
        # Copy events file
        events_dest = base_name + '_events.json'
        shutil.copy(args.events, events_dest)
        
        print(f"Results saved to {args.save}", file=sys.stderr)
        print(f"Metadata saved to {meta_file}", file=sys.stderr)
        print(f"Events saved to {events_dest}", file=sys.stderr)

    if args.output == 'terminal':
        print(ReportFormatter.format_terminal(stats, checks))
    elif args.output == 'csv':
        print(ReportFormatter.format_csv_header())
        print(ReportFormatter.format_csv_row(stats))
    elif args.output == 'json':
        print(ReportFormatter.format_json(stats, checks))

    sys.exit(0 if all(checks.values()) else 1)


if __name__ == '__main__':
    main()
