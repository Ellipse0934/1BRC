package generate

import (
    "fmt"
    "strings"
    "math/rand"
    "io"
    "strconv"
    "1BRC/src/generate/data"
)


// Algorithm for city generation
type NameGen int64

const (
    GenRand NameGen = iota
    GenDefault
    GenExtended
)

func (n NameGen) String() string {
    switch n {
    case GenRand:
        return "Random"
    case GenDefault:
        return "Default"
    case GenExtended:
        return "Extended"
    }
    return "unknown"
}

func ParseGen(s string) (NameGen) {
    switch s {
    case "random":
        return GenRand
    case "default":
        return GenDefault
    case "extended":
        return GenExtended
    default:
        panic("Unknown option")
    }
}

var charsetBasic = []rune("abcdefghijklmnopqrstuvwxyz")

func randstring(r *rand.Rand, length int) string {
    b := strings.Builder{}
    for i := 0; i < length; i++ {
        b.WriteRune(charsetBasic[r.Int() % len(charsetBasic)])
    }
    return b.String()
}

type station struct {
    name string
    meanTemp int
}

// Generate n cities and output n samples for each city
// City name is between 50 and 100 chars, Normal distribution
// City temperature is 
func generate(seed int64, ncities int, nsamples int) {
    r := rand.New(rand.NewSource(seed))
    // name -> samples
    m := make(map[string]bool)
    for i := 0; i < ncities; i++ {
        l := 50 + (r.Int() % 50)
        name := randstring(r, l)
        for _, ok := m[name]; ok ;{
            l = 50 + (r.Int() % 50)
            name = randstring(r, l)
        } 
        fmt.Printf("%s\n", name)
    }
}

func GenerateReal(rows int,w io.Writer) {
    s := make([]byte, 10)
    for i := 0; i < 10; i++ {
        ss := strconv.Itoa(i)
        copy(s, ss)
        w.Write(s)
        w.Write([]byte{'\n'})
    }
}

func randtemp(r *rand.Rand) int {
    var t int = int(r.NormFloat64() * 150 + 270)
    t = max(t, -999)
    t = min(t, 999)
    return t 
}

func randselection(r *rand.Rand, n int, stations []string) []string {
    // TODO: Add assertion n > len(stations)
    // TODO: If n > len(stations) / 2, choose rejections
    m := make(map[int]bool)
    result := make([]string, 0, len(stations))
    for i := 0; i < n; i++ {
        selection := r.Int() % len(stations)
        for _, ok := m[selection]; ok != false; {
            selection += 1
            _, ok = m[selection]
        }
        m[selection] = true
        x := string(stations[selection])
        result = append(result, x)
    }
    return result
}

// Given `NameGen` make a map of `nstations -> mean_temp`
func GenerateStations(r *rand.Rand, nstations int, s NameGen) ([]station) {
    result := make([]station, 0, nstations)
    switch s {
    case GenRand:
        m := make(map[string]bool)
        for i := 0; i < nstations; i++ {
            l := 50 + (r.Int() % 50)
            name := randstring(r, l)
            for _, ok := m[name]; ok != false; {
                l = 20 + (r.Int() % 20)
                name = randstring(r, l)
            }
            m[name] = true
            result = append(result, station{name, randtemp(r)})
        }
    case GenDefault:
        // TODO: use map merge or smthing
        stations := data.GetDefault()
        for _, name := range randselection(r, nstations, stations) {
            result = append(result, station{name, randtemp(r)})
        }
    case GenExtended:
        stations := data.GetExtended()
        for _, name := range randselection(r, nstations, stations) {
            result = append(result, station{name, randtemp(r)})
        }
    default:
        panic("Not implemented")
    }
    return result
}

// Generate stations, temp vals
// Write a random temp b/w mean +- 15
func Generate(nrows int,nstations int, s NameGen, w io.Writer) {
    r := rand.New(rand.NewSource(0))
    stations := GenerateStations(r, nstations, s)
    line := make([]byte, 0, 128)
    for i := 0; i < nrows; i++ {
        line := line[:0]
        st := stations[r.Int() % len(stations)]
        t := st.meanTemp + (r.Int() % 20) - 10
        t = min(max(t, -999), 999)
        line = append(line, st.name...)
        line = append(line, ';')
        line = strconv.AppendInt(line, int64(t / 10), 10)
        line = append(line, '.')
        line = strconv.AppendInt(line, math.Abs(int64(t % 10), 10))
        line = append(line, '\n')
        w.Write(line)
        //fmt.Fprintf(w, "%s;%0.1f\n", st.name, float64(t) / 10.0)
    }
}

