package api

import(
    "github.com/gorilla/context"
    "github.com/gorilla/mux"
    "net/http"
    "strconv"
)

const PlanetTypeArctic = "arctic"
const PlanetTypeDesert = "desert"
const PlanetTypeOceanic = "oceanic"
const PlanetTypeTemperate = "temperate"
const PlanetTypeTropical = "tropical"
const PlanetTypeRocky = "rocky"
const PlanetTypeVolcanic = "volcanic"

type(
	Planet struct {
		TableName struct{} `json:"-" sql:"map__planets"`
	
		Id uint16 `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Population uint `json:"population"`
		SystemId uint16 `json:"-"`
		System *System `json:"system"`
		OrbitId uint16 `json:"-"`
		Orbit *SystemOrbit `json:"orbit"`
		PlayerId uint16 `json:"-"`
		Player *Player `json:"player"`
		Resources []PlanetResource `json:"resources"`
		StorageId uint16 `json:"-"`
		Storage *Storage `json:"storage"`
		SettingsId uint16 `json:"-"`
		Settings *PlanetSettings `json:"settings"`
		Relations []DiplomaticRelation `json:"relations"`
		Buildings []Building `json:"buildings"`
		NbBuildings uint8 `json:"nb_buildings" sql:"-"`
		AvailableBuildings []BuildingPlan `json:"available_buildings" sql:"-"`
	}
	
	PlanetResource struct {
		TableName struct{} `json:"-" sql:"map__planet_resources"`
	
		Name string `json:"name"`
		Density uint8 `json:"density"`
		PlanetId uint16 `json:"-"`
		Planet *Planet `json:"planet"`
	}

	PlanetData struct {
	  Picto string
	  Image string
	  Resources map[string]uint8
	}
	PlanetsData map[string]PlanetData
)

func GetPlanet(w http.ResponseWriter, r *http.Request) {
    player := context.Get(r, "player").(*Player)
    id, _ := strconv.ParseUint(mux.Vars(r)["id"], 10, 16)

    SendJsonResponse(w, 200, getPlayerPlanet(uint16(id), player.Id))
}

func (p *Planet) changeOwner(player *Player) {
	p.Player = player
    p.PlayerId = player.Id
}

func (p *Planet) update() {
    if err := Database.Update(p); err != nil {
        panic(NewException("Planet could not be updated", err))
    }
}

func getPlanet(id uint16) *Planet {
    planet := &Planet{}
    if err := Database.
        Model(planet).
        Column("planet.*", "Player", "Player.Faction", "Settings", "Relations", "Relations.Player", "Relations.Player.Faction", "Relations.Faction", "Resources", "System", "Storage").
        Where("planet.id = ?", id).
        Select(); err != nil {
            return nil
    }
    return planet
}

func (s *System) getPlanets() []Planet {
    planets := make([]Planet, 0)
    if err := Database.
        Model(&planets).
        Column("planet.*", "Orbit", "Player", "Player.Faction").
        Where("planet.system_id = ?", s.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "System not found", err))
    }
    return planets
}

func (p *Player) getPlanets() []Planet {
    planets := make([]Planet, 0)
    if err := Database.
        Model(&planets).
        Column("planet.*", "Player", "Player.Faction", "System", "Resources", "Settings").
        Where("planet.player_id = ?", p.Id).
        Select(); err != nil {
            panic(NewHttpException(404, "Player not found", err))
    }
    return planets
}

func getPlayerPlanet(id uint16, playerId uint16) *Planet {
    planet := &Planet{}
    if err := Database.
        Model(planet).
        Column("planet.*", "Player", "Player.Faction", "Settings", "Relations", "Relations.Player", "Relations.Player.Faction", "Relations.Faction", "Resources", "System", "Storage").
        Where("planet.id = ?", id).
        Select(); err != nil {
            panic(NewHttpException(404, "Planet not found", err))
    }
    if planet.Player != nil && playerId == planet.Player.Id {
        getPlanetOwnerData(planet)
    }
    return planet
}

func getPlanetOwnerData(planet *Planet) {
    planet.Buildings, planet.AvailableBuildings = planet.getBuildings()
    planet.NbBuildings = 7
}