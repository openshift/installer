#!/usr/bin/env python3
"""Validate test documentation metadata and structure.

Usage:
    python3 check.py                          # validate all files
    python3 check.py cases                    # validate only test case files
    python3 check.py plans                    # validate only test plan files
    python3 check.py matrix                   # validate only the feature matrix
    python3 check.py gcd/cases/gcd_sovereign_install.md  # validate a specific file
"""
import re
import sys
from pathlib import Path

DOCS_DIR = Path(__file__).parent

# -- Enumerations -----------------------------------------------------------

AUTOMATION_STATUSES = {
    "Automated (unit test)",
    "Automated (Prow CI)",
    "Manual",
    "TBD",
}

TEST_TYPES = {"Unit", "Integration", "E2E"}

MATRIX_CELL_VALUES = {"AT", "MT", "AT/MT", "PT", "NT", "NA", "Yes", "No"}

MATRIX_CELL_PATTERN = re.compile(
    r"^(?:"
    r"(?:AT|MT|AT/MT|PT|NT|NA)(?:\[\d+\])?"
    r"|Yes(?:\[\d+\])?"
    r"|No"
    r"|"  # empty cell (section header rows)
    r")$"
)

JIRA_LINK_PATTERN = re.compile(
    r"\[([A-Z]+-\d+)\]\(https://redhat\.atlassian\.net/browse/[A-Z]+-\d+\)"
)

# Required metadata fields for test case files (the Metadata table).
CASE_METADATA_FIELDS = {"Feature", "Component", "Test type", "Assignee"}

# Required top-level sections for test plan files.
PLAN_REQUIRED_SECTIONS = [
    "1. Introduction",
    "2. Testing Strategy",
    "3. Test Areas and Test Cases",
    "4. Test Details",
    "5. Risks",
    "6. Exit Criteria",
]


# -- Helpers ----------------------------------------------------------------


class ValidationError:
    def __init__(self, file, line, message):
        self.file = file
        self.line = line
        self.message = message

    def __str__(self):
        rel = self.file
        if self.line:
            return f"{rel}:{self.line}: {self.message}"
        return f"{rel}: {self.message}"


def relative_path(path):
    """Return path relative to DOCS_DIR, or the path itself if outside."""
    try:
        return path.relative_to(DOCS_DIR)
    except ValueError:
        return path


def parse_md_table(lines, start):
    """Parse a markdown table starting at line index `start`.

    Returns a list of rows, where each row is a list of cell strings
    (stripped of leading/trailing whitespace and pipes).  The separator
    row (|---|---|) is skipped.
    """
    rows = []
    i = start
    while i < len(lines):
        line = lines[i].strip()
        if not line.startswith("|"):
            break
        cells = [c.strip() for c in line.split("|")[1:-1]]
        if cells and not all(re.fullmatch(r"-+|:-+|-+:|:-+:", c) for c in cells):
            rows.append(cells)
        i += 1
    return rows


def find_table_after_heading(lines, heading_text):
    """Find the first markdown table that appears after a heading containing
    `heading_text`.  Returns (header_row, data_rows, line_number) or None."""
    for i, line in enumerate(lines):
        stripped = line.strip().lstrip("#").strip()
        if heading_text.lower() in stripped.lower():
            for j in range(i + 1, min(i + 10, len(lines))):
                if lines[j].strip().startswith("|"):
                    rows = parse_md_table(lines, j)
                    if len(rows) >= 1:
                        return rows[0], rows[1:], j + 1
                    break
    return None


def strip_bold(text):
    """Remove markdown bold markers from text."""
    return text.replace("**", "").strip()


# -- File type detection ----------------------------------------------------


def detect_file_type(path, lines):
    """Detect whether a file is a 'case', 'plan', or 'matrix'."""
    name = path.name
    if name == "platform_feature_matrix.md":
        return "matrix"

    first_heading = ""
    for line in lines[:5]:
        if line.startswith("# "):
            first_heading = line.strip()
            break

    if first_heading.startswith("# Test Case:"):
        return "case"
    if first_heading.startswith("# Feature Test Plan:"):
        return "plan"
    if first_heading.startswith("# Platform Feature Matrix"):
        return "matrix"

    # Fall back to path-based detection.
    parts = path.parts
    if "cases" in parts:
        return "case"
    if "plans" in parts:
        return "plan"

    return "unknown"


# -- Validators -------------------------------------------------------------


def validate_case(path, lines):
    """Validate a test case file."""
    errors = []
    rel = relative_path(path)

    # 1. Title must start with "# Test Case:"
    title_found = False
    for i, line in enumerate(lines):
        if line.startswith("# "):
            if not line.strip().startswith("# Test Case:"):
                errors.append(ValidationError(
                    rel, i + 1,
                    f"title must start with '# Test Case:' (got: {line.strip()!r})",
                ))
            else:
                title_text = line.strip().removeprefix("# Test Case:").strip()
                if not title_text:
                    errors.append(ValidationError(
                        rel, i + 1, "title after '# Test Case:' is empty",
                    ))
            title_found = True
            break
    if not title_found:
        errors.append(ValidationError(rel, None, "no top-level heading found"))

    # 2. Metadata table
    result = find_table_after_heading(lines, "Metadata")
    if result is None:
        errors.append(ValidationError(rel, None, "missing '## Metadata' section with table"))
    else:
        header, rows, line_num = result
        found_fields = {}
        for row_idx, row in enumerate(rows):
            if len(row) < 2:
                errors.append(ValidationError(
                    rel, line_num + row_idx + 2,
                    f"metadata row has fewer than 2 columns: {row}",
                ))
                continue
            field_name = strip_bold(row[0])
            field_value = row[1].strip()
            found_fields[field_name] = (field_value, line_num + row_idx + 2)

        for required in CASE_METADATA_FIELDS:
            if required not in found_fields:
                errors.append(ValidationError(
                    rel, None,
                    f"missing required metadata field: '{required}'",
                ))

        # Feature must contain a JIRA link.
        if "Feature" in found_fields:
            val, ln = found_fields["Feature"]
            if not JIRA_LINK_PATTERN.search(val):
                errors.append(ValidationError(
                    rel, ln,
                    "Feature field must contain a JIRA link "
                    "(e.g. [KEY-123](https://redhat.atlassian.net/browse/KEY-123))",
                ))

        # Test type must be from the allowed set (comma-separated).
        if "Test type" in found_fields:
            val, ln = found_fields["Test type"]
            parts = [p.strip() for p in val.split(",")]
            for part in parts:
                if part not in TEST_TYPES:
                    errors.append(ValidationError(
                        rel, ln,
                        f"invalid test type '{part}' "
                        f"(expected one of: {', '.join(sorted(TEST_TYPES))})",
                    ))

        # Assignee must not be empty.
        if "Assignee" in found_fields:
            val, ln = found_fields["Assignee"]
            if not val:
                errors.append(ValidationError(
                    rel, ln, "Assignee field is empty",
                ))

    # 3. Scenarios table
    table_statuses = {}  # scenario number -> status from table
    result = find_table_after_heading(lines, "Scenarios")
    if result is None:
        errors.append(ValidationError(rel, None, "missing '## Scenarios' section with table"))
    else:
        header, rows, line_num = result
        if len(header) < 3:
            errors.append(ValidationError(
                rel, line_num,
                f"scenarios table must have at least 3 columns (got {len(header)})",
            ))
        has_manual = False
        for row_idx, row in enumerate(rows):
            if len(row) < 3:
                continue
            status = row[2].strip()
            if status not in AUTOMATION_STATUSES:
                errors.append(ValidationError(
                    rel, line_num + row_idx + 2,
                    f"invalid automation status '{status}' "
                    f"(expected one of: {', '.join(sorted(AUTOMATION_STATUSES))})",
                ))
            if status == "Manual":
                has_manual = True

            num_str = row[0].strip()
            if not num_str.isdigit() or int(num_str) < 1:
                errors.append(ValidationError(
                    rel, line_num + row_idx + 2,
                    f"scenario number must be a positive integer (got '{num_str}')",
                ))
            else:
                table_statuses[int(num_str)] = status

        if has_manual:
            full_text = "\n".join(lines)
            if "MANUAL_TESTS_START" not in full_text:
                errors.append(ValidationError(
                    rel, None,
                    "file has Manual scenarios but missing "
                    "<!-- MANUAL_TESTS_START --> marker",
                ))
            if "MANUAL_TESTS_END" not in full_text:
                errors.append(ValidationError(
                    rel, None,
                    "file has Manual scenarios but missing "
                    "<!-- MANUAL_TESTS_END --> marker",
                ))

    # 4. Per-scenario automation status - cross-check with scenarios table.
    detail_statuses = {}  # scenario number -> (status, line)
    scenario_pattern = re.compile(r"^##\s+(\d+)\.\s+")
    for i, line in enumerate(lines):
        m = scenario_pattern.match(line)
        if m:
            num = int(m.group(1))
            found_status = False
            for j in range(i + 1, min(i + 8, len(lines))):
                if "**Automation status:**" in lines[j]:
                    found_status = True
                    status_val = lines[j].split("**Automation status:**")[1].strip()
                    if status_val not in AUTOMATION_STATUSES:
                        errors.append(ValidationError(
                            rel, j + 1,
                            f"invalid per-scenario automation status '{status_val}'",
                        ))
                    detail_statuses[num] = (status_val, j + 1)
                    break
            if not found_status:
                errors.append(ValidationError(
                    rel, i + 1,
                    f"scenario {num} missing '**Automation status:**' field",
                ))

    # Cross-check: statuses in the table must match per-scenario statuses.
    for num, table_status in table_statuses.items():
        if num in detail_statuses:
            detail_status, detail_line = detail_statuses[num]
            if table_status != detail_status:
                errors.append(ValidationError(
                    rel, detail_line,
                    f"scenario {num} status mismatch: table says "
                    f"'{table_status}' but detail says '{detail_status}'",
                ))
        elif num not in detail_statuses:
            errors.append(ValidationError(
                rel, None,
                f"scenario {num} listed in table but has no "
                f"'## {num}. ...' detail section",
            ))

    # 5. Related Links section should exist.
    if not any(line.strip().startswith("## Related Links") for line in lines):
        errors.append(ValidationError(
            rel, None, "missing '## Related Links' section",
        ))

    return errors


def validate_plan(path, lines):
    """Validate a test plan file."""
    errors = []
    rel = relative_path(path)

    # 1. Title must start with "# Feature Test Plan:"
    title_found = False
    for i, line in enumerate(lines):
        if line.startswith("# "):
            if not line.strip().startswith("# Feature Test Plan:"):
                errors.append(ValidationError(
                    rel, i + 1,
                    f"title must start with '# Feature Test Plan:' "
                    f"(got: {line.strip()!r})",
                ))
            else:
                title_text = line.strip().removeprefix("# Feature Test Plan:").strip()
                if not title_text:
                    errors.append(ValidationError(
                        rel, i + 1, "title after '# Feature Test Plan:' is empty",
                    ))
            title_found = True
            break
    if not title_found:
        errors.append(ValidationError(rel, None, "no top-level heading found"))

    # 2. Required sections.
    section_headings = []
    for i, line in enumerate(lines):
        m = re.match(r"^##\s+(.+)", line)
        if m:
            section_headings.append((m.group(1).strip(), i + 1))

    found_sections = {h for h, _ in section_headings}
    for required in PLAN_REQUIRED_SECTIONS:
        if not any(required.lower() in h.lower() for h in found_sections):
            errors.append(ValidationError(
                rel, None, f"missing required section: '## {required}'",
            ))

    # 3. Introduction must contain Overview, Scope, Key Features, References.
    intro_subsections = {"Overview", "Scope", "Key Features", "References"}
    found_subsections = set()
    for i, line in enumerate(lines):
        m = re.match(r"^###\s+\d+\.\d+\.\s+(.+)", line)
        if m:
            found_subsections.add(m.group(1).strip())
    for sub in intro_subsections:
        if sub not in found_subsections:
            errors.append(ValidationError(
                rel, None,
                f"Introduction missing subsection '### x.x. {sub}'",
            ))

    # 4. Must contain at least one JIRA link.
    full_text = "\n".join(lines)
    if not JIRA_LINK_PATTERN.search(full_text):
        errors.append(ValidationError(
            rel, None,
            "test plan must contain at least one JIRA link "
            "(e.g. [KEY-123](https://redhat.atlassian.net/browse/KEY-123))",
        ))

    # 5. Testing Strategy should list test types.
    if not any("Unit Tests" in line or "Unit tests" in line for line in lines):
        errors.append(ValidationError(
            rel, None,
            "Testing Strategy should describe unit test coverage",
        ))

    # 6. Test Cases table should exist in section 3.2, with valid case links.
    result = find_table_after_heading(lines, "Test Cases")
    if result is None:
        errors.append(ValidationError(
            rel, None,
            "section 3 should contain a 'Test Cases' table linking to case docs",
        ))
    else:
        _, case_rows, case_line_num = result
        link_pattern = re.compile(r"\[([^\]]+)\]\(([^)]+)\)")
        for row_idx, row in enumerate(case_rows):
            for cell in row:
                for match in link_pattern.finditer(cell):
                    link_target = match.group(2)
                    if link_target.startswith("http://") or link_target.startswith("https://"):
                        continue
                    if link_target.startswith("/") or ".." not in link_target:
                        pass  # relative path, check it resolves
                    resolved = (path.parent / link_target).resolve()
                    if not resolved.exists():
                        errors.append(ValidationError(
                            rel, case_line_num + row_idx + 2,
                            f"case link target does not exist: '{link_target}'",
                        ))

    # 7. Exit Criteria should have at least one bullet.
    in_exit = False
    has_criteria = False
    for line in lines:
        if re.match(r"^##\s+.*Exit Criteria", line, re.IGNORECASE):
            in_exit = True
            continue
        if in_exit:
            if line.startswith("## "):
                break
            if line.strip().startswith("- "):
                has_criteria = True
                break
    if not has_criteria:
        errors.append(ValidationError(
            rel, None,
            "Exit Criteria section should have at least one bullet point",
        ))

    return errors


def find_section_table(lines, section_heading):
    """Find a table under a heading that exactly matches `## <section_heading>`.

    Unlike find_table_after_heading, this matches the exact section heading
    (case-insensitive) and only searches within that section (stops at the
    next ## heading).
    """
    in_section = False
    for i, line in enumerate(lines):
        if re.match(r"^##\s+" + re.escape(section_heading) + r"\s*$", line, re.IGNORECASE):
            in_section = True
            continue
        if in_section:
            if line.startswith("## "):
                break
            if line.strip().startswith("|"):
                rows = parse_md_table(lines, i)
                if rows:
                    return rows[0], rows[1:], i + 1
    return None


def validate_matrix(path, lines):
    """Validate the platform feature matrix file."""
    errors = []
    rel = relative_path(path)

    # 1. Find the Legend table and extract valid codes.
    legend_result = find_section_table(lines, "Legend")
    if legend_result is None:
        errors.append(ValidationError(rel, None, "missing '## Legend' section with table"))
    else:
        _, legend_rows, _ = legend_result
        legend_codes = set()
        for row in legend_rows:
            if row:
                legend_codes.add(row[0].strip())
        expected_codes = {"AT", "MT", "AT/MT", "PT", "NT", "NA"}
        missing = expected_codes - legend_codes
        if missing:
            errors.append(ValidationError(
                rel, None,
                f"Legend table missing codes: {', '.join(sorted(missing))}",
            ))

    # 2. Find the Feature Matrix table.
    result = find_section_table(lines, "Feature Matrix")
    if result is None:
        errors.append(ValidationError(
            rel, None, "missing '## Feature Matrix' section with table",
        ))
        return errors

    header, rows, line_num = result

    # Header should have Feature + at least one platform.
    if len(header) < 2:
        errors.append(ValidationError(
            rel, line_num,
            f"matrix table header must have at least 2 columns (got {len(header)})",
        ))
        return errors

    platform_count = len(header) - 1  # first column is "Feature"

    # 3. Validate each row.
    referenced_footnotes = set()
    footnote_ref_pattern = re.compile(r"\[(\d+)\]")

    for row_idx, row in enumerate(rows):
        row_line = line_num + row_idx + 2  # +2: 1-based + separator row
        feature = row[0].strip() if row else ""

        # Section header rows have bold text and empty platform cells.
        is_section_header = feature.startswith("**") and feature.endswith("**")

        if len(row) != len(header):
            errors.append(ValidationError(
                rel, row_line,
                f"row has {len(row)} columns but header has {len(header)}: "
                f"'{feature}'",
            ))
            continue

        for col_idx in range(1, len(row)):
            cell = row[col_idx].strip()

            # Collect footnote references.
            for m in footnote_ref_pattern.finditer(cell):
                referenced_footnotes.add(int(m.group(1)))

            # Strip footnote references for validation.
            clean_cell = footnote_ref_pattern.sub("", cell).strip()

            if not MATRIX_CELL_PATTERN.match(cell):
                if is_section_header and clean_cell == "":
                    continue
                errors.append(ValidationError(
                    rel, row_line,
                    f"invalid cell value '{cell}' in column "
                    f"'{header[col_idx]}' for feature '{feature}' "
                    f"(expected one of: {', '.join(sorted(MATRIX_CELL_VALUES))} "
                    f"with optional [N] footnote)",
                ))

    # 4. Validate footnotes section exists and all referenced footnotes are
    #    defined.
    defined_footnotes = set()
    in_footnotes = False
    footnote_def_pattern = re.compile(r"^(\d+)\.\s+")
    for i, line in enumerate(lines):
        if re.match(r"^##\s+Footnotes\s*$", line):
            in_footnotes = True
            continue
        if in_footnotes:
            if line.startswith("## "):
                break
            m = footnote_def_pattern.match(line.strip())
            if m:
                defined_footnotes.add(int(m.group(1)))

    if not in_footnotes:
        errors.append(ValidationError(rel, None, "missing '## Footnotes' section"))

    undefined = referenced_footnotes - defined_footnotes
    if undefined:
        errors.append(ValidationError(
            rel, None,
            f"footnotes referenced but not defined: "
            f"{', '.join(str(n) for n in sorted(undefined))}",
        ))

    unused = defined_footnotes - referenced_footnotes
    if unused:
        errors.append(ValidationError(
            rel, None,
            f"footnotes defined but not referenced in table: "
            f"{', '.join(str(n) for n in sorted(unused))}",
        ))

    # 5. Platform Key table should exist.
    platform_key = find_section_table(lines, "Platform Key")
    if platform_key is None:
        errors.append(ValidationError(rel, None, "missing '## Platform Key' section"))

    return errors


# -- Discovery and dispatch -------------------------------------------------


def find_files(filter_arg=None):
    """Find markdown files to validate based on CLI argument."""
    if filter_arg is None:
        # All files: matrix + all cases + all plans.
        files = []
        matrix = DOCS_DIR / "platform_feature_matrix.md"
        if matrix.exists():
            files.append(matrix)
        for md in sorted(DOCS_DIR.rglob("*.md")):
            if md == matrix or md.name == "README.md":
                continue
            files.append(md)
        return files

    if filter_arg == "matrix":
        matrix = DOCS_DIR / "platform_feature_matrix.md"
        return [matrix] if matrix.exists() else []

    if filter_arg == "cases":
        return sorted(DOCS_DIR.rglob("*/cases/*.md"))

    if filter_arg == "plans":
        return sorted(DOCS_DIR.rglob("*/plans/*.md"))

    # Specific path (relative to DOCS_DIR or absolute).
    target = Path(filter_arg)
    if not target.is_absolute():
        target = DOCS_DIR / target
    target = target.resolve()

    docs_root = DOCS_DIR.resolve()
    if not str(target).startswith(str(docs_root) + "/") and target != docs_root:
        print(
            f"error: path is outside docs directory: {filter_arg}",
            file=sys.stderr,
        )
        sys.exit(2)

    if target.exists():
        return [target]

    print(f"error: file not found: {filter_arg}", file=sys.stderr)
    sys.exit(2)


def validate_file(path):
    """Validate a single file, returning a list of ValidationErrors."""
    lines = path.read_text().splitlines()
    file_type = detect_file_type(path, lines)

    if file_type == "case":
        return validate_case(path, lines)
    elif file_type == "plan":
        return validate_plan(path, lines)
    elif file_type == "matrix":
        return validate_matrix(path, lines)
    else:
        return [ValidationError(
            relative_path(path), None,
            "could not detect file type (expected case, plan, or matrix)",
        )]


def main():
    filter_arg = sys.argv[1] if len(sys.argv) > 1 else None

    if filter_arg in ("-h", "--help"):
        print(__doc__.strip())
        sys.exit(0)

    files = find_files(filter_arg)
    if not files:
        print("No files to validate.", file=sys.stderr)
        sys.exit(2)

    total_errors = 0
    for path in files:
        errors = validate_file(path)
        for err in errors:
            print(f"  {err}")
        if errors:
            total_errors += len(errors)
        else:
            rel = relative_path(path)
            print(f"  {rel}: ok")

    print()
    files_label = "file" if len(files) == 1 else "files"
    if total_errors:
        print(
            f"FAIL: {total_errors} error(s) in {len(files)} {files_label}",
        )
        sys.exit(1)
    else:
        print(f"OK: {len(files)} {files_label} validated")
        sys.exit(0)


if __name__ == "__main__":
    main()
