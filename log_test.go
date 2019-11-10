package log

import (
	"testing"
	"github.com/rs/zerolog/log"
)
/**/
func BenchmarkWriteToFile(b *testing.B) {
	log.Logger = log.With().Caller().Logger()
	Init()
	SetLogLevel(0)
	// DisableStdOut()
	// EnableStdOut()
	for i := 0; i < b.N; i++ {
		log.Debug().Msg("testing log msg jhjjj")
	}
}
func BenchmarkWriteToFile5Filed(b *testing.B) {
	log.Logger = log.With().Caller().Logger()
	InitWithPathAndFormat("../", "2006_01_02")
	SetLogLevel(1)
	// DisableStdOut()
	// EnableStdOut()
	for i := 0; i < b.N; i++ {
		log.Info().Int16("int16", 16).Float32("float32", 32.32).Str("string", "string").Int64("int64", 64).Msg("testing filed")
	}
}

type xxx struct {

}
func (x *xxx) Write(p []byte) (n int, e error) {
	return 0, nil
}

func BenchmarkNoWrite(b *testing.B) {
	log.Logger = log.Output(&xxx{})
	for i := 0; i < b.N; i++ {
		log.Debug().Msg("testing log msg jhjjj")
	}
}
func BenchmarkNoWrite5Filed(b *testing.B) {
	log.Logger = log.Output(&xxx{})
	for i := 0; i < b.N; i++ {
		log.Info().Int16("int16", 16).Float32("float32", 32.32).Str("string", "string").Int64("int64", 64).Msg("testing filed")
	}
}