#!/usr/bin/env python3
"""Validate test documentation metadata and structure.

Usage:
    python3 check.py                          # validate all files
    python3 check.py cases                    # validate only test case files
    python3 check.py plans                    # validate only test plan files
    python3 check.py matrix                   # validate only the feature matrix
    python3 check.py gcd/cases/gcd_private_dns_only.md   # validate a specific file
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

PRIORITY_VALUES = {"P1", "P2", "P3", ""}

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

# Required metadata fields for test case files (the table under the title).
CASE_METADATA_FIELDS = ["Feature", "Component", "Type", "Priority", "Test plan"]

# Required sections for test case files, in order (IEEE-829 shape).
CASE_REQUIRED_SECTIONS = ["Setup", "Test", "Cleanup"]

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


def find_leading_table(lines):
    """Find the first markdown table that appears before any '## ' heading.

    Returns (header_row, data_rows, line_number) or None.
    """
    for i, line in enumerate(lines):
        if line.startswith("## "):
            return None
        if line.strip().startswith("|"):
            rows = parse_md_table(lines, i)
            if rows:
                return rows[0], rows[1:], i + 1
    return None


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
    """Validate a test case file.

    A test case file is one manual (or not-yet-automated) scenario written in
    the IEEE-829 shape: a metadata table, then Setup / Test / Cleanup, with
    the Test section made of alternating Step and Expect subsections.
    """
    errors = []
    rel = relative_path(path)

    # 1. Title must start with "# Test Case:".
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

    # 2. Metadata table, directly under the title.
    result = find_leading_table(lines)
    if result is None:
        errors.append(ValidationError(
            rel, None,
            "missing metadata table between the title and the first section",
        ))
    else:
        _, rows, line_num = result
        found_fields = {}
        for row_idx, row in enumerate(rows):
            if len(row) < 2:
                errors.append(ValidationError(
                    rel, line_num + row_idx + 2,
                    f"metadata row has fewer than 2 columns: {row}",
                ))
                continue
            found_fields[strip_bold(row[0])] = (
                row[1].strip(), line_num + row_idx + 2,
            )

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

        # Type must be a known automation status.
        if "Type" in found_fields:
            val, ln = found_fields["Type"]
            if val not in AUTOMATION_STATUSES:
                errors.append(ValidationError(
                    rel, ln,
                    f"invalid type '{val}' "
                    f"(expected one of: {', '.join(sorted(AUTOMATION_STATUSES))})",
                ))

        # Priority must be P1/P2/P3 or empty.
        if "Priority" in found_fields:
            val, ln = found_fields["Priority"]
            if val not in PRIORITY_VALUES:
                errors.append(ValidationError(
                    rel, ln,
                    f"invalid priority '{val}' (expected one of: P1, P2, P3, or empty)",
                ))

        # Test plan must link to a file that exists.
        if "Test plan" in found_fields:
            val, ln = found_fields["Test plan"]
            m = re.search(r"\[[^\]]+\]\(([^)]+)\)", val)
            if not m:
                errors.append(ValidationError(
                    rel, ln,
                    "Test plan field must be a markdown link to the parent plan",
                ))
            elif not (path.parent / m.group(1)).resolve().exists():
                errors.append(ValidationError(
                    rel, ln,
                    f"test plan link target does not exist: '{m.group(1)}'",
                ))

    # 3. Setup / Test / Cleanup sections, in order.
    sections = [
        (line.strip().removeprefix("## ").strip(), i + 1)
        for i, line in enumerate(lines)
        if line.startswith("## ")
    ]
    section_names = [name for name, _ in sections]
    for required in CASE_REQUIRED_SECTIONS:
        if required not in section_names:
            errors.append(ValidationError(
                rel, None, f"missing required section: '## {required}'",
            ))
    present = [n for n in section_names if n in CASE_REQUIRED_SECTIONS]
    expected_order = [n for n in CASE_REQUIRED_SECTIONS if n in present]
    if present != expected_order:
        errors.append(ValidationError(
            rel, None,
            f"sections must appear in the order "
            f"{' -> '.join(CASE_REQUIRED_SECTIONS)} (got: {' -> '.join(present)})",
        ))

    # 4. The Test section must alternate '### Step' and '### Expect'.
    in_test = False
    steps = []
    for i, line in enumerate(lines):
        if line.startswith("## "):
            in_test = line.strip() == "## Test"
            continue
        if in_test and line.startswith("### "):
            steps.append((line.strip().removeprefix("### ").strip(), i + 1))

    if not any(name == "Step" for name, _ in steps):
        errors.append(ValidationError(
            rel, None, "the Test section must contain at least one '### Step'",
        ))
    for idx, (name, ln) in enumerate(steps):
        expected = "Step" if idx % 2 == 0 else "Expect"
        if name != expected:
            errors.append(ValidationError(
                rel, ln,
                f"Test subsections must alternate Step/Expect: "
                f"expected '### {expected}', got '### {name}'",
            ))
            break
    else:
        if len(steps) % 2 != 0:
            errors.append(ValidationError(
                rel, steps[-1][1],
                "the last '### Step' has no matching '### Expect'",
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
        return None


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
    validated = 0
    skipped = 0
    for path in files:
        errors = validate_file(path)
        if errors is None:
            rel = relative_path(path)
            print(f"  {rel}: skipped (unknown file type)")
            skipped += 1
            continue
        validated += 1
        for err in errors:
            print(f"  {err}")
        if errors:
            total_errors += len(errors)
        else:
            rel = relative_path(path)
            print(f"  {rel}: ok")

    print()
    files_label = "file" if validated == 1 else "files"
    summary = f"OK: {validated} {files_label} validated"
    if skipped:
        summary += f", {skipped} skipped"
    if total_errors:
        print(
            f"FAIL: {total_errors} error(s) in {validated} {files_label}",
        )
        sys.exit(1)
    else:
        print(summary)
        sys.exit(0)


if __name__ == "__main__":
    main()
