# Fuzzy Finder – Rules & Information Design Document

## 1. Purpose

This document defines the **rules, data, and algorithms** required to implement a high‑performance fuzzy finder suitable for:

* File search
* Command palettes
* Symbol navigation
* General text lookup

The goal is to balance **match quality**, **ranking correctness**, and **real‑time performance**.

---

## 2. Core Definitions

### 2.1 Query

A short string entered by the user (typically 1–50 characters).

### 2.2 Candidate

A searchable string from a collection (e.g., filename, command name, symbol).

### 2.3 Match

A query matches a candidate if the query is a **subsequence** of the candidate.

### 2.4 Score

A numeric value representing match quality. Higher scores rank higher.

---

## 3. Matching Rules

### 3.1 Subsequence Rule (Mandatory)

Characters in the query must appear in the candidate **in order**, but not necessarily contiguously.

Example:

* Query: `gitc`
* Candidate: `getInitialConfig` → MATCH

### 3.2 Case Sensitivity Rule

Configurable behavior:

* Case‑insensitive by default
* Optional case‑sensitive bonus

### 3.3 Normalization Rule

Before matching:

* Convert to lowercase (unless case‑sensitive mode)
* Normalize Unicode (NFKD recommended)
* Strip or normalize diacritics (optional)

---

## 4. Scoring System

### 4.1 Base Score

Each matched character contributes a base score.

Example:

* +10 points per matched character

### 4.2 Positional Bonuses

| Rule               | Description             | Typical Weight |
| ------------------ | ----------------------- | -------------- |
| Consecutive Match  | Adjacent characters     | +5 to +15      |
| Word Start         | Start of string or word | +10            |
| Separator Boundary | After `/ - _ space .`   | +8             |
| CamelCase Boundary | lower → Upper           | +8             |
| Exact Case Match   | `A` → `A`               | +2             |

### 4.3 Positional Penalties

| Rule           | Description            | Typical Weight |
| -------------- | ---------------------- | -------------- |
| Gap Penalty    | Characters skipped     | −1 per char    |
| Late Match     | Match occurs late      | −1 per index   |
| Long Candidate | Prefer shorter strings | −length factor |

### 4.4 Score Formula (Conceptual)

```
score = Σ(base + bonuses − penalties)
```

---

## 5. Ranking Rules

### 5.1 Primary Sort

Sort candidates by descending score.

### 5.2 Tie‑Breakers (in order)

1. Fewer gaps
2. Shorter candidate length
3. Earlier first match index
4. Lexicographical order

---

## 6. Algorithms

### 6.1 Matching Algorithm Options

#### Option A: Greedy Subsequence Scan (Recommended Default)

* Single forward pass
* Fast (O(n))
* Good quality with smart scoring

#### Option B: Dynamic Programming (Optional)

* Similar to LCS / edit distance
* Better theoretical scoring
* Slower (O(n × m))

#### Option C: Bitmask Matching (Advanced)

* Uses bitwise operations
* Extremely fast for short queries (<64 chars)

---

## 7. Required Data Per Candidate

Each candidate should precompute:

```text
- Original string
- Lowercased string
- Character positions
- Word boundary indexes
- Separator indexes
- CamelCase transitions
- Length
```

This data enables fast scoring without recomputation.

---

## 8. Performance Strategies

### 8.1 Incremental Search

When query grows:

* Filter previous results instead of full dataset

### 8.2 Early Exit Rules

Abort matching when:

* Remaining characters < remaining query length
* Maximum possible score < current threshold

### 8.3 Result Limiting

Keep only top‑K results (e.g., 50–200) using:

* Min‑heap
* Partial sort

### 8.4 Caching

Cache:

* Query prefix results
* Candidate preprocessing

---

## 9. Configuration Options

| Option           | Description                 |
| ---------------- | --------------------------- |
| Case sensitivity | On / Off / Smart            |
| Scoring weights  | Tunable per rule            |
| Max results      | UI dependent                |
| Ignore patterns  | `.git`, `node_modules`, etc |

---

## 10. UI Integration Requirements

### 10.1 Real‑Time Constraints

* Target <16ms per keystroke
* Non‑blocking UI updates

### 10.2 Highlighting Data

Matcher should return:

* Match indexes
* Grouped ranges for rendering

---

## 11. Failure & Edge Cases

* Empty query → return top recent / popular items
* No matches → return empty state
* Very long candidates → early penalty
* Unicode combining characters

---

## 12. Testing Requirements

### 12.1 Correctness Tests

* Subsequence validation
* Boundary detection

### 12.2 Ranking Tests

* Known query → expected order

### 12.3 Performance Tests

* Large dataset (100k+ entries)
* Worst‑case queries

---

## 13. Success Criteria

A fuzzy finder is considered successful if it:

* Feels instantaneous
* Ranks intuitive results first
* Degrades gracefully with dataset size
* Is configurable without algorithm changes

---

## 14. Non‑Goals

* Full spell correction
* Semantic meaning matching
* NLP‑based intent inference

---

## 15. Summary

A fuzzy finder is fundamentally:

> **Subsequence matching + scoring rules + aggressive optimization**

This document defines the minimum rules and information needed to implement one that feels modern, fast, and intuitive.
