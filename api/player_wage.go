package api

import(
    "math"
    "sync"
)

func CalculatePlayersWage() {
    nbPlayers, _ := Database.Model(&Player{}).Count()
    limit := 20

    var wg sync.WaitGroup

    for offset := 0; offset < nbPlayers; offset +=limit {
        players := getPlayers(offset, limit)

        for _, player := range players {
          if player.IsActive==true {
            wg.Add(1)
            go player.calculateWage(&wg)
          }
        }
        wg.Wait()
    }
}

func (p *Player) calculateWage(wg *sync.WaitGroup) {
    defer wg.Done()
    defer CatchException()
    baseWage := int32(50)
    serviceWageRatio := float64(0.5)
    wage := int32(0)
    for _, planet := range p.getPlanets() {
      wage += baseWage +  int32( math.Round( float64(planet.Settings.ServicesPoints) * serviceWageRatio))
    }
    p.updateWallet(wage)
    p.update()
}

func getPlayers(offset int, limit int) []Player {
    players := make([]Player, 0)
    if err := Database.
        Model(&players).
        Limit(limit).
        Offset(offset).
        Select(); err != nil {
            panic(NewException("Players could not be retrieved", err))
    }
    return players
}
