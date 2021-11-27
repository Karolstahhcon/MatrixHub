#!/usr/bin/env bash

# Copyright 2026 MatrixHub Authors
# SPDX-License-Identifier: Apache-2.0

GENERATED_START='<!-- BEGIN GENERATED RELEASE NOTES -->'
GENERATED_END='<!-- END GENERATED RELEASE NOTES -->'
CHANGELOG_WORK_DIR=''

usage() {
  cat <<'EOF'
Usage:
  scripts/changelog.sh --version vX.Y.Z [options]

Options:
  --repo OWNER/REPO       GitHub repository (defaults to GITHUB_REPOSITORY)
  --start-ref REF         Previous official release (auto-detected when omitted)
  --end-ref REF           Candidate commit or branch (default: HEAD)
  --base-branch BRANCH    PR base branch used to disambiguate associated PRs
  --output PATH           CHANGELOG file to update (default: CHANGELOG/CHANGELOG-X.Y.md)
  --dry-run               Print the generated version section without writing a file
  --help                  Show this help
EOF
}

git_cmd() {
  git "$@"
}

github_api() {
  local endpoint=$1
  shift

  gh api \
    -H 'Accept: application/vnd.github+json' \
    -H 'X-GitHub-Api-Version: 2022-11-28' \
    "$@" \
    "$endpoint"
}

http_status_code() {
  awk 'NR == 1 && /^HTTP\// { print $2; exit }'
}

http_response_body() {
  awk '
    NR == 1 && /^HTTP\// { headers = 1; next }
    headers {
      sub(/\r$/, "")
      if ($0 == "") headers = 0
      next
    }
    { print }
  '
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "error: required command not found: $1" >&2
    return 1
  }
}

append_line() {
  printf '%s\n' "$2" >> "$1"
}

has_line() {
  local lines=$1
  local expected=$2
  local line

  while IFS= read -r line; do
    [ "$line" = "$expected" ] && return 0
  done <<EOF
$lines
EOF
  return 1
}

trim_blank_lines() {
  awk '
    { lines[NR] = $0 }
    END {
      first = 1
      while (first <= NR && lines[first] ~ /^[[:space:]]*$/) first++
      last = NR
      while (last >= first && lines[last] ~ /^[[:space:]]*$/) last--
      if (first <= last) {
        sub(/^[[:space:]]+/, "", lines[first])
        sub(/[[:space:]]+$/, "", lines[last])
      }
      for (i = first; i <= last; i++) print lines[i]
    }
  '
}

extract_release_note() {
  awk '
    function remember(value) {
      content[++count] = value
    }
    {
      if (done) next
      line = $0
      sub(/\r$/, "", line)
      lower = tolower(line)
      if (!found && match(lower, /```release-note[[:blank:]]*/)) {
        found = 1
        rest = substr(line, RSTART + RLENGTH)
        closing = index(rest, "```")
        if (closing) {
          remember(substr(rest, 1, closing - 1))
          closed = 1
          done = 1
          next
        }
        remember(rest)
        next
      }
      if (found) {
        closing = index(line, "```")
        if (closing) {
          remember(substr(line, 1, closing - 1))
          closed = 1
          done = 1
          next
        }
        remember(line)
      }
    }
    END {
      if (closed) {
        for (i = 1; i <= count; i++) print content[i]
      }
    }
  ' | trim_blank_lines
}

is_no_release_note() {
  local upper
  upper=$(printf '%s' "$1" | tr '[:lower:]' '[:upper:]')
  [[ $upper =~ ^[^A-Z0-9_]*(NONE|NO)[^A-Z0-9_]*$ ]]
}

version_is_less() {
  local left=${1#v}
  local right=${2#v}
  local -a left_parts right_parts
  local index left_part right_part left_digit right_digit

  IFS=. read -r -a left_parts <<EOF
$left
EOF
  IFS=. read -r -a right_parts <<EOF
$right
EOF

  for index in 0 1 2; do
    left_part=${left_parts[$index]}
    right_part=${right_parts[$index]}
    while [ "${#left_part}" -gt 1 ] && [ "${left_part#0}" != "$left_part" ]; do
      left_part=${left_part#0}
    done
    while [ "${#right_part}" -gt 1 ] && [ "${right_part#0}" != "$right_part" ]; do
      right_part=${right_part#0}
    done
    [ "${#left_part}" -lt "${#right_part}" ] && return 0
    [ "${#left_part}" -gt "${#right_part}" ] && return 1
    while [ -n "$left_part" ]; do
      left_digit=${left_part%"${left_part#?}"}
      right_digit=${right_part%"${right_part#?}"}
      [ "$left_digit" -lt "$right_digit" ] && return 0
      [ "$left_digit" -gt "$right_digit" ] && return 1
      left_part=${left_part#?}
      right_part=${right_part#?}
    done
  done
  return 1
}

find_previous_official_tag() {
  local version=$1
  local end_ref=$2
  local tag tags tag_commit first_parent_commits
  local best=''

  tags=$(git_cmd tag --merged "$end_ref" --list) || return 1
  first_parent_commits=$(git_cmd rev-list --first-parent "$end_ref") || return 1
  while IFS= read -r tag; do
    [[ $tag =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || continue
    tag_commit=$(git_cmd rev-parse "$tag^{commit}") || return 1
    has_line "$first_parent_commits" "$tag_commit" || continue
    if [ -z "$best" ] || version_is_less "$best" "$tag"; then
      best=$tag
    fi
  done <<EOF
$tags
EOF

  if [ -n "$best" ] && ! version_is_less "$best" "$version"; then
    echo "error: --version $version must be greater than latest official tag $best on $end_ref" >&2
    return 1
  fi
