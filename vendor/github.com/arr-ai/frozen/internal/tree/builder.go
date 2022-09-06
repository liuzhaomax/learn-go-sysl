package tree

import (
	"fmt"

	"github.com/arr-ai/frozen/internal/depth"
)

var (
	DefaultNPKeyEqArgs      = NewDefaultKeyEqArgs(depth.NonParallel)
	DefaultNPKeyCombineArgs = NewCombineArgs(DefaultNPKeyEqArgs, UseRHS)
)

func NewDefaultKeyEqArgs(gauge depth.Gauge) *EqArgs {
	return NewEqArgs(gauge, elementEqual, hashValue, hashValue)
}

// Builder provides a more efficient way to build nodes incrementally.
type Builder struct {
	t Tree
}

func NewBuilder(capacity int) *Builder {
	return &Builder{}
}

func (b *Builder) Count() int {
	return b.t.count
}

func (b *Builder) Add(args *CombineArgs, v elementT) {
	if b.t.root == nil {
		b.t.root = newLeaf1(v)
		b.t.count = 1
	} else {
		h := newHasher(v, 0)
		if vetting {
			backup := b.clone()
			defer vet(func() { backup.Add(args, v) }, &b.t)(nil)
		}
		var matches int
		b.t.root, matches = b.t.root.Add(args, v, 0, h)
		b.t.count += 1 - matches
	}
}

func (b *Builder) Remove(args *EqArgs, v elementT) {
	if b.t.root != nil {
		h := newHasher(v, 0)
		if vetting {
			backup := b.clone()
			defer vet(func() { backup.Remove(args, v) }, &b.t)(nil)
		}
		var matches int
		b.t.root, matches = b.t.root.Remove(args, v, 0, h)
		b.t.count -= matches
	}
}

func (b *Builder) Get(args *EqArgs, el elementT) *elementT {
	return b.t.Get(args, el)
}

func (b *Builder) Finish() Tree {
	t := b.Borrow()
	*b = Builder{}
	return t
}

func (b *Builder) Borrow() Tree {
	return b.t
}

func (b Builder) String() string {
	return b.Borrow().String()
}

func (b Builder) Format(state fmt.State, verb rune) {
	b.Borrow().Format(state, verb)
}

func (b *Builder) clone() *Builder {
	return &Builder{t: b.t.clone()}
}
