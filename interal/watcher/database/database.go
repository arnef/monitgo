package database

import "github.com/arnef/monitgo/pkg"

type Database interface {
	OnSnapshot(snap []pkg.NodeSnapshot)
}
