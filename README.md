# matpw: Matrix Password Manager Go

[toc]

## Features

* easy to generate password
* Even if matrix is leaked, if permutation of password flagment is not leaked, it is safe
* generate matrix that contains password fragments
* simple and intuitive CLI UX

## Installation

download binary from [Github Releases](https:/example.com)

## What is Matrix Password

```txt
   ---------------------
   | a | b | c | d | e |
   | f | g | h | i | j |
   | k | l | m | n | o |
   | p | q | r | s | t |
   | u | v | w | x | y |
   ----------------------
```

There is 5x5 matrix like the above.
First, you determine the cells four times as permutation. (deplication is not allowd)
for example, select a->e->u->y. select c->w->t->p

Well here, if you determine pattern as a->e->u->y.
Then, generate matrix password, and if result is below, your password is `6E{+<sc5`

```txt
   --------------------------
   | 6E | Sy | x1 | Aw | {+ |
   | !3 | $6 | Ui | 64 | 64 |
   | _T | 6] | 6] | #< | #< |
   | +C | +C | (b | %q | R& |
   | <s | K2 | c5 | Yx | c5 |
   --------------------------
```

## Quiq start

```sh
matpw create
```

first, input service from prompt (required)

```txt
input service:
```

second, input account from prompt (required)

```txt
input account:
```

finally, input descripiton (option)

```txt
input descripiton:
```

then, password matrix will be generated as below

```txt
   --------------------------
   | 6E | Sy | x1 | Aw | {+ |

   | !3 | $6 | Ui | 64 | 64 |

   | _T | 6] | 6] | #< | #< |

   | +C | +C | (b | %q | R& |

   | <s | K2 | c5 | Yx | c5 |
   --------------------------
```

## Usage

### Create

`matpw create`

input`service title` and `account` and `description`. Then password can be generated.

### Search

`matpw search`

You can search password from enrolled service title incrementally.

## LICENSE

Apache 2.0
