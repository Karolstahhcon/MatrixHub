#!/usr/bin/env bash

# Copyright 2026 MatrixHub Authors
# SPDX-License-Identifier: Apache-2.0

set -u
set -o pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck source=changelog.sh
source "$SCRIPT_DIR/changelog.sh"

TEST_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/matrixhub-changelog-test.XXXXXX")
PASSED=0
FAILED=0

cleanup() {
  if [ -n "${TEST_ROOT:-}" ] && [ -d "$TEST_ROOT" ]; then
    rm -rf "$TEST_ROOT"
  fi
}
trap cleanup EXIT HUP INT TERM

fail() {
  printf '  %s\n' "$1" >&2
  return 1
}

assert_equal() {
  [ "$1" = "$2" ] || fail "expected <$1>, got <$2>"
}

assert_contains() {
  case "$1" in
    *"$2"*) return 0 ;;
    *) fail "expected output to contain <$2>" ;;
  esac
}

assert_not_contains() {
  case "$1" in
    *"$2"*) fail "expected output not to contain <$2>" ;;
    *) return 0 ;;
  esac
}

new_work_dir() {
  local directory
  directory=$(mktemp -d "$TEST_ROOT/work.XXXXXX") || return 1
  mkdir -p "$directory/entries"
  : > "$directory/errors"
  : > "$directory/warnings"
  : > "$directory/excluded"
  printf '%s\n' "$directory"
}

pull_json() {
  local number=$1
  local body=$2
  local labels=$3
  local author=${4:-contributor}

  jq -cn \
    --argjson number "$number" \
    --arg body "$body" \
    --argjson labels "$labels" \
    --arg author "$author" \
    '{
      number: $number,
      html_url: ("https://github.com/matrixhub-ai/matrixhub/pull/" + ($number | tostring)),
      body: $body,
      labels: ($labels | map({name: .})),
      user: {login: $author}
    }'
}

test_release_note_parsing() {
  local actual
  actual=$(printf 'before\r\n```release-note\r\nFirst line.\r\nSecond line.\r\n```\r\nafter\n' | extract_release_note)
  assert_equal $'First line.\nSecond line.' "$actual" || return 1
  actual=$(printf 'prefix ```release-note Inline note.``` suffix\n' | extract_release_note)
  assert_equal 'Inline note.' "$actual" || return 1
  actual=$(printf '```release-note\n  Padded note.  \n```\n' | extract_release_note)
  assert_equal 'Padded note.' "$actual" || return 1
  actual=$({
    printf '```release-note\nLarge body note.\n```\n'
    awk 'BEGIN { for (i = 0; i < 20000; i++) print "body suffix" }'
  } | extract_release_note) || return 1
  assert_equal 'Large body note.' "$actual" || return 1
  assert_equal '' "$(printf '```release-note\nunclosed\n' | extract_release_note)" || return 1
  is_no_release_note ' NONE ' || return 1
  is_no_release_note 'no.' || return 1
  if is_no_release_note '_NONE_'; then
    fail 'underscores are word characters in the ci-bot contract'
  fi
  if is_no_release_note 'No user-visible change'; then
    fail 'descriptive text must not be treated as NONE'
  fi
}

test_version_and_kind_rules() {
  version_is_less v0.9.9 v0.10.0 || return 1
  if version_is_less v0.10.0 v0.9.9; then
    fail 'numeric version ordering is reversed'
  fi
  version_is_less v99999999999999999999.0.0 v100000000000000000000.0.0 || return 1
  assert_equal api-change "$(classify_kind $'kind/feature\nkind/api-change')" || return 1
  assert_equal bug-or-regression "$(classify_kind 'kind/regression')" || return 1
  assert_equal other "$(classify_kind 'kind/cleanup')" || return 1
}

test_previous_official_tag() (
  local include_v010=false

  git_cmd() {
    case "$1" in
      tag)
        printf '%s\n' v0.1.0 v0.2.0-rc.1 v0.2.0 v0.9.0 v0.10.0 not-a-version
        ;;
      rev-list)
        [ "$2" = --first-parent ] || return 90
        printf '%s\n' head
        [ "$include_v010" = false ] || printf '%s\n' commit-0.10.0
        printf '%s\n' commit-0.2.0 commit-0.1.0
        ;;
      rev-parse)
        case "$2" in
          'v0.1.0^{commit}') printf '%s\n' commit-0.1.0 ;;
          'v0.2.0^{commit}') printf '%s\n' commit-0.2.0 ;;
          'v0.9.0^{commit}') printf '%s\n' side-commit-0.9.0 ;;
          'v0.10.0^{commit}') printf '%s\n' commit-0.10.0 ;;
          *) return 91 ;;
        esac
        ;;
      *) return 92 ;;
    esac
  }

  assert_equal v0.2.0 "$(find_previous_official_tag v0.3.0 HEAD)" || return 1
  assert_equal v0.2.0 "$(find_previous_official_tag v0.10.0 HEAD)" || return 1
  include_v010=true
  assert_equal v0.10.0 "$(find_previous_official_tag v0.11.0 HEAD)" || return 1
  if find_previous_official_tag v0.3.0 HEAD >/dev/null 2>&1; then
    fail 'a target version older than the latest official tag must be rejected'
  fi
  is_first_parent_ref v0.2.0 HEAD || return 1
  if is_first_parent_ref v0.9.0 HEAD; then
    fail 'a tag reachable only through a side branch must not be a release boundary'
  fi
)

test_pull_validation_and_rendering() {
  local work feature none label_only_none empty_none content_with_none dependabot_body block errors
  work=$(new_work_dir) || return 1
  feature=$(pull_json 21 $'```release-note\nAdded a capability.\nMore detail.\n```' \
    '["kind/feature"]' alice) || return 1
  none=$(pull_json 22 $'```release-note\nNONE\n```' \
    '["kind/cleanup", "release-note-none"]' bob) || return 1
  dependabot_body=$'Bumps a dependency from 1.0.0 to 1.0.1.\n\n<details>\n<summary>Release notes</summary>\nUpstream notes.\n</details>'
  label_only_none=$(pull_json 23 "$dependabot_body" \
    '["kind/dependency", "release-note-none"]' 'dependabot[bot]') || return 1
  empty_none=$(pull_json 24 $'```release-note\n\n```' \
    '["kind/cleanup", "release-note-none"]' bob) || return 1
  content_with_none=$(pull_json 25 $'```release-note\nBlock content wins.\n```' \
    '["kind/bug", "release-note", "release-note-none"]' carol) || return 1

  process_pull "$feature" "$work" || return 1
  process_pull "$none" "$work" || return 1
  process_pull "$label_only_none" "$work" || return 1
  process_pull "$empty_none" "$work" || return 1
  process_pull "$content_with_none" "$work" || return 1
