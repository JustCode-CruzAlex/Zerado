// Package cli is the command-line surface treated as what it is: a public API
// that gets designed, versioned and kept stable, rather than accreted.
//
// A CLI is an API the moment anyone types it. A verb name in somebody's shell
// history, a script parsing an exit code, a cron line — each is a caller with
// no upgrade path, and getting a verb name wrong is a breaking change later.
// So the surface is declared here, in one table, with its exit codes and its
// output shape, and the stability policy is written down before there is
// anything to break.
//
// # What this package is and is not
//
// It is the surface: verbs, arguments, flags, exit codes, output envelope,
// phase and stability. It contains no flag parsing, no command implementations
// and no printing. The point is that the surface can be reviewed, diffed and
// tested as data before any of it is built, and that a future command cannot
// quietly acquire a verb that was never designed.
package cli

// APIVersion is the major version of the machine-readable output contract.
//
// It appears in every --json envelope. A script may branch on it. It is
// incremented only for a breaking change to the envelope or to a documented
// field's meaning; adding a field is not breaking and does not bump it.
const APIVersion = 1

// Phase is when a verb becomes real.
//
// It is on the surface rather than in a comment because of anti-pattern 14:
// copy that mentions a capability the build does not have presents something
// unbuilt as working. A reserved verb is listed in the documentation as
// reserved and refuses with a clear message; it never half-works.
type Phase uint8

const (
	// PhaseOne verbs ship in Phase 1.
	PhaseOne Phase = iota

	// PhaseLater verbs are reserved: the name is claimed so a Phase 1 user
	// cannot squat it and so a future release does not have to pick a worse
	// word. They are not routed and they are not advertised in help as
	// available.
	PhaseLater
)

// Stability is the promise attached to a verb.
type Stability uint8

const (
	// StabilityStable means the verb, its arguments, its exit codes and its
	// JSON fields will not change incompatibly within this major version.
	StabilityStable Stability = iota

	// StabilityReserved means the name is claimed and nothing else is
	// promised.
	StabilityReserved
)

// Verb is one top-level command.
type Verb struct {
	// Name is what the player types. It is the API.
	Name string

	// SummaryKey names the catalogue entry describing it in help. A key
	// rather than a sentence: help text is user-facing, and D9 admits no
	// literal even here.
	SummaryKey string

	// Args are the positional arguments, in order.
	Args []Arg

	// Flags are the options this verb accepts, on top of [GlobalFlags].
	Flags []Flag

	// Phase and Stability say whether this verb does anything yet and what is
	// promised about it.
	Phase     Phase
	Stability Stability

	// Interactive reports whether the verb hands control to the TUI rather
	// than printing and exiting.
	//
	// It matters to the contract because an interactive verb cannot honour
	// --json, and a caller piping one deserves a usage error rather than a
	// terminal full of escape sequences.
	Interactive bool

	// NeedsNetwork records the verb's offline class, from
	// 07-offline-contract.md's three-way split. It is here so the CLI's own
	// help can state it and so a scripted caller can tell, before running
	// anything, which verbs will refuse on a plane.
	NeedsNetwork bool
}

// Arg is a positional argument.
type Arg struct {
	Name string

	// Required means the verb refuses without it — unless one of
	// [Arg.RequiredUnless] is present.
	Required bool

	// RequiredUnless names flags whose presence makes this argument optional.
	//
	// It exists because the surface has to be able to express `mark --clear`,
	// and an earlier revision could not: `state` was flatly Required and
	// `--clear` was declared alongside it, so a parser reading this surface as
	// data would reject `zerado mark <game> --clear` for a missing argument —
	// while the dossier says plainly that a CLI without --clear could not do
	// what Z-06 does.
	//
	// A surface that contradicts its own stated requirement is worse than an
	// incomplete one, because the contradiction is only discovered by whoever
	// writes the parser. [Verb.MissingArgs] resolves it, and
	// TestTheSurfaceCanExpressEveryDocumentedInvocation asserts the surface
	// admits each invocation the dossier promises.
	RequiredUnless []string

	// Repeated means the argument may be given more than once and collects.
	Repeated bool
}

// Flag is one option.
type Flag struct {
	// Name is the long form, without dashes. There are no single-letter
	// aliases in Phase 1: a one-letter flag is the hardest thing to rename
	// and the easiest to collide, and the TUI is where speed lives.
	Name string

	// SummaryKey names the catalogue entry describing it.
	SummaryKey string

	// Value is the placeholder shown in help, or empty for a boolean flag.
	Value string
}

// GlobalFlags apply to every verb.
//
// They are global because each one is about how the program talks rather than
// about what it does, and because a flag that means the same thing everywhere
// should not have to be re-declared per verb — which is how it comes to mean
// three slightly different things.
func GlobalFlags() []Flag {
	return []Flag{
		{Name: "json", SummaryKey: "cli.flag.json"},
		{Name: "quiet", SummaryKey: "cli.flag.quiet"},
		{Name: "db", SummaryKey: "cli.flag.db", Value: "path"},
		{Name: "no-color", SummaryKey: "cli.flag.no_color"},
		{Name: "version", SummaryKey: "cli.flag.version"},
		{Name: "help", SummaryKey: "cli.flag.help"},
	}
}

// Surface returns every verb Zerado claims, in help order.
//
// The reserved entries are as much a part of the deliverable as the live ones:
// they are the names Phase 2 and Phase 3 will need, claimed now, so that a
// Phase 1 user's script cannot come to depend on `zerado tonight` meaning
// something else and so that the eventual feature does not have to settle for
// a worse word.
func Surface() []Verb {
	return []Verb{
		{
			// The bare command is the product. Everything else is a
			// convenience for people who are not in it.
			Name:        "",
			SummaryKey:  "cli.verb.tui",
			Phase:       PhaseOne,
			Stability:   StabilityStable,
			Interactive: true,
		},
		{
			Name:       "sync",
			SummaryKey: "cli.verb.sync",
			Args:       []Arg{{Name: "provider"}},
			Flags: []Flag{
				{Name: "all", SummaryKey: "cli.flag.sync.all"},
			},
			Phase:        PhaseOne,
			Stability:    StabilityStable,
			NeedsNetwork: true,
		},
		{
			Name:       "list",
			SummaryKey: "cli.verb.list",
			Flags: []Flag{
				{Name: "state", SummaryKey: "cli.flag.list.state", Value: "state"},
				{Name: "source", SummaryKey: "cli.flag.list.source", Value: "provider"},
				{Name: "search", SummaryKey: "cli.flag.list.search", Value: "text"},
				{Name: "absent", SummaryKey: "cli.flag.list.absent"},
				{Name: "limit", SummaryKey: "cli.flag.list.limit", Value: "n"},
			},
			Phase:     PhaseOne,
			Stability: StabilityStable,
		},
		{
			Name:       "mark",
			SummaryKey: "cli.verb.mark",
			Args: []Arg{
				{Name: "game", Required: true},
				{Name: "state", Required: true, RequiredUnless: []string{"clear"}},
			},
			Flags: []Flag{
				// Clearing an override is a different action from setting
				// NOT STARTED, and the CLI has to be able to express both or
				// it cannot do what Z-06 does. That is why `state` above is
				// required *unless* this flag is present rather than flatly
				// required.
				{Name: "clear", SummaryKey: "cli.flag.mark.clear"},
			},
			Phase:     PhaseOne,
			Stability: StabilityStable,
		},
		{
			Name:       "add",
			SummaryKey: "cli.verb.add",
			Flags: []Flag{
				{Name: "title", SummaryKey: "cli.flag.add.title", Value: "text"},
				{Name: "platform", SummaryKey: "cli.flag.add.platform", Value: "text"},
				{Name: "state", SummaryKey: "cli.flag.add.state", Value: "state"},
				{Name: "owned-since", SummaryKey: "cli.flag.add.owned_since", Value: "date"},
			},
			Phase:     PhaseOne,
			Stability: StabilityStable,
		},
		{
			// The deep link. 04-navigation-and-focus.md §288 notes that this
			// requires the route stack to be constructible from a descriptor
			// rather than only by pushing — which is a Phase 1 constraint that
			// exists because of this verb, and is recorded here so the cost is
			// attached to the thing that causes it.
			Name:        "game",
			SummaryKey:  "cli.verb.game",
			Args:        []Arg{{Name: "game", Required: true}},
			Phase:       PhaseOne,
			Stability:   StabilityStable,
			Interactive: true,
		},
		{
			// Everything Z-09 reports, for someone who cannot open a TUI —
			// the database path, the vault backing, the image capability, the
			// audio state. It reaches no network, which is what makes it a
			// diagnostic rather than a connectivity test.
			Name:       "doctor",
			SummaryKey: "cli.verb.doctor",
			Phase:      PhaseOne,
			Stability:  StabilityStable,
		},
		{
			Name:       "version",
			SummaryKey: "cli.verb.version",
			Phase:      PhaseOne,
			Stability:  StabilityStable,
		},
		{
			Name:       "help",
			SummaryKey: "cli.verb.help",
			Args:       []Arg{{Name: "verb"}},
			Phase:      PhaseOne,
			Stability:  StabilityStable,
		},

		// ---- Reserved. Claimed, not routed, never half-working. ----
		{Name: "tonight", SummaryKey: "cli.verb.tonight", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "price", SummaryKey: "cli.verb.price", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "watch", SummaryKey: "cli.verb.watch", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "tag", SummaryKey: "cli.verb.tag", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "enrich", SummaryKey: "cli.verb.enrich", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "devices", SummaryKey: "cli.verb.devices", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "export", SummaryKey: "cli.verb.export", Phase: PhaseLater, Stability: StabilityReserved},
		{Name: "import", SummaryKey: "cli.verb.import", Phase: PhaseLater, Stability: StabilityReserved},
	}
}

// MissingArgs returns the names of this verb's required positional arguments
// that were not supplied, given the flags that were.
//
// given is the set of positional arguments actually present, in order, and
// flags is the set of flag names supplied. It is the one piece of behaviour on
// this surface, and it is here rather than in a future parser because the
// surface's own claim — that it can express every documented invocation — is
// otherwise unfalsifiable.
func (v Verb) MissingArgs(given []string, flags map[string]bool) []string {
	var missing []string
	for i, a := range v.Args {
		if i < len(given) && given[i] != "" {
			continue
		}
		if !a.Required {
			continue
		}
		if satisfiedBy(a.RequiredUnless, flags) {
			continue
		}
		missing = append(missing, a.Name)
	}
	return missing
}

// satisfiedBy reports whether any of names is present in flags.
func satisfiedBy(names []string, flags map[string]bool) bool {
	for _, n := range names {
		if flags[n] {
			return true
		}
	}
	return false
}

// Lookup returns the verb with this name.
func Lookup(name string) (Verb, bool) {
	for _, v := range Surface() {
		if v.Name == name {
			return v, true
		}
	}
	return Verb{}, false
}
