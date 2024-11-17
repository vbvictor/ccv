#!/usr/bin/env python3
import subprocess
import argparse
from collections import defaultdict
import re
import json

def get_most_modified_files(args):
	if args.exclude:
		exclude_regex = re.compile(args.exclude)

	if args.ext:
		allowed_extensions = set('.' + ext.strip() for ext in args.ext.split(','))

	if args.verbose:
		print(f"Starting git log analysis{'for the last ' + str(args.top) + ' commits' if args.top else ' for all commits'}...")
		if args.ext:
			print(f"Including only files with extensions: {', '.join(allowed_extensions)}")
		if args.since and args.until:
			print(f"Analyzing commits between {args.since} and {args.until}")
		elif args.since:
			print(f"Analyzing commits since {args.since}")
		elif args.until:
			print(f"Analyzing commits until {args.until}")

	cmd = ['git', 'log', '--pretty=format:%H', '--numstat']

	if args.top:
		cmd.insert(2, f'-n{args.top}')

	if args.since:
		cmd.append(f'--since={args.since}')

	if args.until:
		cmd.append(f'--until={args.until}')

	if args.path:
		cmd.append('--')
		cmd.append(args.path)

	if args.verbose:
		print("Executing git log command...", cmd)
	result = subprocess.run(cmd, capture_output=True, encoding='utf-8', errors='replace')

	if args.verbose:
		print("Processing git log data...")
	file_changes = defaultdict(lambda: {
		'changes': 0, 
		'additions': 0,
		'deletions': 0,
		'commits': 0
	})

	current_commit = None
	modified_in_commit = set()
	commits_processed = 0
	total_commits = args.top if args.top else len([l for l in result.stdout.splitlines() if len(l) == 40])

	for line in result.stdout.splitlines():
		if not line.strip():
			continue
	
		if len(line) == 40:  # This is a commit hash
			if current_commit and modified_in_commit:
				for filepath in modified_in_commit:
					if args.exclude and exclude_regex.search(filepath):
						continue
					if args.ext and not any(filepath.endswith(ext) for ext in allowed_extensions):
						continue
					file_changes[filepath]['commits'] += 1
			current_commit = line
			modified_in_commit = set()
			commits_processed += 1
			if args.verbose and commits_processed % 100 == 0:
				print(f"Processing commit {commits_processed}/{total_commits} ({(commits_processed/total_commits)*100:.1f}%)")
		else:
			parts = line.split()
			if len(parts) == 3 and parts[0].isdigit() and parts[1].isdigit():
				additions, deletions, filepath = parts
				additions = int(additions)
				deletions = int(deletions)
				if args.exclude and exclude_regex.search(filepath):
					continue
				if args.ext and not any(filepath.endswith(ext) for ext in allowed_extensions):
					continue
				file_changes[filepath]['additions'] += additions
				file_changes[filepath]['deletions'] += deletions
				file_changes[filepath]['changes'] += additions + deletions
				modified_in_commit.add(filepath)

	if args.verbose:
		print(f"Sorting results by {args.sort}...")
	
	sorted_files = sorted(file_changes.items(), key=lambda x: x[1][args.sort], reverse=True)
	
	if args.verbose:
		print(f"Analysis complete. Found {len(file_changes)} unique files.")
	return sorted_files[:args.top]

def print_tabular_results(most_modified, args):
	sort_description = {
		'changes': 'total lines changed',
		'additions': 'lines added',
		'deletions': 'lines deleted',
		'commits': 'commit count'
	}[args.sort]

	print(f"\nTop {args.top} most modified files (by {sort_description}):")
	print("-" * 100)
	print(f"{'CHANGES':>8} {'ADDED':>8} {'DELETED':>8} {'COMMITS':>8} {'FILEPATH'}")
	print("-" * 100)
	for filepath, info in most_modified:
		print(f"{info['changes']:8d} {info['additions']:8d} {info['deletions']:8d} {info['commits']:8d} {filepath}")

def print_json_results(most_modified, args):
	output = {
		"metadata": {
			"total_files": len(most_modified),
			"args.sort": args.sort,
			"filters": {
				"path": args.path,
				"exclude_pattern": args.exclude,
				"extensions": args.ext,
				"date_range": {
					"since": args.since,
					"until": args.until
				}
			}
		},
		"files": [
			{
				"path": filepath,
				"changes": info["changes"],
				"additions": info["additions"],
				"deletions": info["deletions"],
				"commits": info["commits"]
			}
			for filepath, info in most_modified
		]
	}
	print(json.dumps(output, indent=2))

def main():

	parser = argparse.ArgumentParser(description="Find the most modified files in a Git repository.")
	parser.add_argument('--commits', '-n', type=int, help="Optional: Number of commits to analyze from the git log", default=None)
	parser.add_argument('--sort', '-s', 
						choices=['changes', 'additions', 'deletions', 'commits'], 
						default='changes',
						help="Sort by: total changes, additions, deletions or commit count (default: changes)")
	parser.add_argument('--top', '-t', type=int, default=10,
						help="Number of top files to display (default: 10)")
	parser.add_argument('--args.verbose', '-v', action='store_true',
						help="Show detailed progress information")
	parser.add_argument('--path', '-p', type=str, default=None,
						help="Only analyze files within specified path")
	parser.add_argument('--exclude', '-e', type=str, default=None,
						help="Exclude files matching this regex pattern")
	parser.add_argument('--ext', type=str, default=None,
						help="Only include files with specified extensions (comma-separated, e.g. 'h,hpp,c,cpp')")
	parser.add_argument('--since', type=str, default=None,
						help="Start date for analysis (e.g. '2023-01-01' or '2 weeks ago')")
	parser.add_argument('--until', type=str, default=None,
						help="End date for analysis (e.g. '2023-12-31' or 'yesterday')")
	parser.add_argument('--format', '-f', choices=['table', 'json'], default='table',
						help="Output format (default: table)")
	args = parser.parse_args()

	most_modified = get_most_modified_files(args)
	if args.format == 'json':
		print_json_results(most_modified, args)
	else:
		print_tabular_results(most_modified, args)

if __name__ == "__main__":
	main()
