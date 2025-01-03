# parc

parc is a [Parser Combinator](https://en.wikipedia.org/wiki/Parser_combinator) written in Golang.

It mainly inspired by the [Arcsecond](https://github.com/francisrstokes/arcsecond) library.
Many ideas are also borrowed from the 
[How Parser Combinators Work, Parser Combinators From Scratch](https://www.youtube.com/watch?v=6oQLRhw5Ah0&t=185s)
youtube series of Low Byte Productions,
the presentation of Jeremy Brown at GopherCon 2022: on _Parsing w/o Code Generators: Parser Combinators in Go with Generics_,
and the [nom](https://docs.rs/nom/latest/nom/index.html) Rust parser combinator library.

See the [tutorial/](tutorial/) to get started with this package.

Read the [parc package doc pages](https://pkg.go.dev/github.com/tombenke/parc).

As every parser has its own unit tests, so find them in the source code,
and study the ones you want to use, in order to better understanding how to use them.

## References

- [Arcsecond](https://github.com/francisrstokes/arcsecond):
  a zero-dependency, Fantasy Land compliant JavaScript Parser Combinator library largely inspired by Haskell's Parsec.
- [You could have invented Parser Combinators](https://theorangeduck.com/page/you-could-have-invented-parser-combinators)
- [How Parser Combinators Work, Parser Combinators From Scratch](https://www.youtube.com/watch?v=6oQLRhw5Ah0&t=185s)
- [Parser-Combinators-From-Scratch](https://github.com/lowbyteproductions/Parser-Combinators-From-Scratch)
- [GopherCon 2022: Jeremy Brown - Parsing w/o Code Generators: Parser Combinators in Go with Generics](https://www.youtube.com/watch?v=x5p_SJNRB4U)
- [Simple parser combinator package as shown at GopherCon 2022](https://github.com/jhbrown-veradept/gophercon22-parser-combinators/tree/main)
- [nom](https://docs.rs/nom/latest/nom/index.html)

