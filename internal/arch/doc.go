// Package arch holds no code. It holds the tests that keep the seam
// architecture true over time.
//
// Three of Zerado's published promises are not properties of any one package —
// they are properties of who imports whom:
//
//   - "Runs with the network off" holds because there is exactly one kind of
//     package that may originate network I/O.
//   - "No telemetry running in the background" is provable by inspection for
//     the same reason.
//   - "A screen never talks to a provider" holds because the import graph does
//     not allow it.
//
// 07-offline-contract.md §7.3 proposes to keep the first two with a grep-level
// rule in review. A rule nobody can check is a convention, and conventions
// lose to deadlines — so the rules are executable here instead, and a change
// that breaks one fails CI with the reason attached rather than being noticed
// by whoever happened to review that afternoon.
package arch
