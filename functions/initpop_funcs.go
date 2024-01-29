package helper

import (
	// "fmt"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Function to create initial population
func CreateInitialPopulation(size int, chromosomes []Chromosome, fas []Player, free_positions map[int][]string, week string, streamable_players []Player) {

	// Create random seed
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)

	// Create WaitGroup
	var wg sync.WaitGroup
	ch := make(chan Chromosome)

	// Create (size) goroutines to generate chromosomes concurrently
	for i := 0; i < size; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chromosome := CreateChromosome(streamable_players, week, fas, free_positions, rng)
			ch <- chromosome
		}()
	}

	// Wait for all goroutines to finish and collect the chromosomes
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collect the chromosomes from the channel
	i := 0
	for chromosome := range ch {
		chromosomes[i] = chromosome
		i++
	}

}

func CreateChromosome(streamable_players []Player, week string, fas []Player, free_positions map[int][]string, rng *rand.Rand) Chromosome {

	streamable_count := len(streamable_players)
	cur_streamers := make([]Player, streamable_count)
	days_in_week := ScheduleMap[week].GameSpan

	chromosome := Chromosome{Genes: make([]Gene, days_in_week+1), FitnessScore: 0, TotalAcquisitions: 0, CumProbTracker: 0.0, DroppedPlayers: make(map[string]DroppedPlayer)}

	// Initialize genes
	for j := 0; j <= days_in_week; j++ {
		chromosome.Genes[j] = Gene{Roster: make(map[string]Player), NewPlayers: make(map[string]Player), Day: j, Acquisitions: 0, DroppedPlayers: []Player{}, Bench: []Player{}}
	}

	// Insert streamable players into chromosome
	InsertStreamablePlayers(streamable_players, free_positions, week, &chromosome, cur_streamers)

	// Create copy of free agents and free positions
	fas_copy := make([]Player, len(fas))
	copy(fas_copy, fas)

	free_positions_copy := make(map[int][]string)
	for k, v := range free_positions {
		free_positions_copy[k] = make([]string, len(v))
		copy(free_positions_copy[k], v)
	}

	// Fill the gene for each day
	for j := 0; j <= days_in_week; j++ {

		// Get gene for day
		gene := &chromosome.Genes[j]

		acq_count := rng.Intn(3 + rng.Intn(2))

		// Check if there are enough available positions to make acq_count moves
		if len(free_positions_copy[j]) < acq_count {
			acq_count = len(free_positions_copy[j])
		}

		// On the first day, make it so you cant drop streamable players who are playing
		if j == 0 && acq_count > len(streamable_players) - len(gene.Roster) {
			acq_count = len(streamable_players) - len(gene.Roster)
		}

		// Add acq_count players to gene
		for k := 0; k < acq_count; k++ {

			// Get random free agent that fits into free_positions
			trials := 0
			cont := true
			for cont && trials < 100 {

				// Get random free agent
				rand_index := rand.Intn(len(fas_copy))
				fa := fas_copy[rand_index]

				// Check if the free agent has a game on the day
				if !Contains(ScheduleMap[week].Games[fa.Team], j) || fa.Injured {
					trials++
					continue
				}

				// Check if the free agent has a valid position on the day
				has_match := false
				for _, pos := range fa.ValidPositions {
					if Contains(free_positions_copy[j], pos) {
						has_match = true
						break
					}
				}

				if has_match {

					add := false

					// Choose the positon that results in the highest "net rostering". Adjusted for choosing the best position for each day
					pos_map := GetPosMap(fa, 
										&chromosome, 
										free_positions_copy, 
										j, 
										week,
										cur_streamers,
										streamable_players,
										j == 0,
										true,
										&add,
										)

					// Remove position from free_positions_copy and player from free agents. Only remove from free pos on inital day so it doesn't interfere with same day moves
					fas_copy = Remove(fas_copy, rand_index)
					
					// When added here, counts as a new player
					if _, ok := pos_map[j]; ok && add {
						gene.NewPlayers[fa.Name] = fa
					} else {
						continue
					}

					// Fill other game days with added player because on other days, he can go on the bench
					for day, pos := range pos_map {

						// Add player to gene for that day. When added here, doesn't count as a new player
						chromosome.Genes[day].Roster[pos] = fa
					}
					// Once a player is added, add another player or go to next day
					break
				}
				trials++
			}
			trials = 0
		}

		// After each day, decrement countdown for dropped players
		for name, player := range chromosome.DroppedPlayers {
			player.Countdown--
			if player.Countdown == 0 {
				delete(chromosome.DroppedPlayers, name)
				// Add player back to free agents
				fas_copy = append(fas_copy, player.Player)
			} else {
				chromosome.DroppedPlayers[name] = player
			}
		}

		// After each day, add the length of NewPlayers to Acquisitions
		gene.Acquisitions += len(gene.NewPlayers)
	}

	return chromosome
}


// Function to find the best position to put a player in
func GetPosMap(player Player, 
			  chromosome *Chromosome,
		      free_positions map[int][]string, 
			  start_day int, 
			  week string,
			  cur_streamers []Player,
			  streamable_players []Player,
			  first_day bool,
			  not_initial bool,
			  add *bool,
			  ) map[int]string {

	pos_map := make(map[int]string)
	if first_day {
	// If it is the first day, either put intial streamers in or adjust for other than initial streamers

		FindBestPositions(player, chromosome, free_positions, pos_map, start_day, week, add)

		// After adding the initial streamers, when doing regular pickups on the first day, drop the streamer with the lowest average points
		if not_initial {

			// Delete the worst bench streamer
			sort.Slice(chromosome.Genes[0].Bench, func(i, j int) bool {
				return chromosome.Genes[0].Bench[i].AvgPoints < chromosome.Genes[0].Bench[j].AvgPoints 
			})

			// Drop the worst player from the chromosome
			DeleteAllOccurrences(chromosome, cur_streamers, player, chromosome.Genes[0].Bench[0], week, start_day)
		}


	} else if len(chromosome.Genes[start_day].Bench) > 0 {
	// If there are streamers on the bench, find best position for new player and drop the worst bench streamer

		// PrintPopulation(*chromosome, free_positions)
		// fmt.Println(start_day)
		// fmt.Println("Not playing streamers:", not_playing_streamers)

		// Get the worst player that is not playing on the day
		sort.Slice(chromosome.Genes[start_day].Bench, func(i, j int) bool {
			return chromosome.Genes[start_day].Bench[i].AvgPoints < chromosome.Genes[start_day].Bench[j].AvgPoints })

		if len(chromosome.Genes[start_day].Bench) == 0 {
			fmt.Println("Not playing streamers:", chromosome.Genes[start_day].Bench)
			fmt.Println("Streamers:", cur_streamers)
		}
		worst_player := chromosome.Genes[start_day].Bench[0]

		// Drop the worst player from the chromosome
		DeleteAllOccurrences(chromosome, cur_streamers, player, worst_player, week, start_day)
		FindBestPositions(player, chromosome, free_positions, pos_map, start_day, week, add)

	} else {
	// If there are no streamers on the bench (ie. the roster is full), drop the worst playing streamer that the fa can fit into and find best position

		sort.Slice(cur_streamers, func(i, j int) bool {
			return cur_streamers[i].AvgPoints < cur_streamers[j].AvgPoints
		})

		i := -1
		for j := 0; j < len(cur_streamers); j++ {
			position := MapContainsValue(chromosome.Genes[start_day].Roster, cur_streamers[j].Name)
			if Contains(ScheduleMap[week].Games[player.Team], start_day) && Contains(player.ValidPositions, position) {
				i = j
				break
			}
		}

		if i == -1 {
			return pos_map
		}
		worst_player := cur_streamers[i]

		DeleteAllOccurrences(chromosome, cur_streamers, player, worst_player, week, start_day)
		// Remove player from new players
		delete(chromosome.Genes[start_day].NewPlayers, worst_player.Name)
		FindBestPositions(player, chromosome, free_positions, pos_map, start_day, week, add)

	}

	return pos_map
}


// Function to insert initial streamable players into a chromosome
func InsertStreamablePlayers(streamable_players []Player, free_positions map[int][]string, week string, chromosome *Chromosome, cur_streamers []Player) {

	sort.Slice(streamable_players, func(i, j int) bool {
		return streamable_players[i].AvgPoints > streamable_players[j].AvgPoints
	})

	// Insert streamable players into chromosome
	for i, player := range streamable_players {

		// Choose the positon that results in the highest "net rostering". Adjusted for choosing the best position for each day
		add := false
		pos_map := GetPosMap(player, 
							chromosome, 
							free_positions, 
							0, 
							week, 
							cur_streamers,
							streamable_players,
							true,
							false,
							&add,
							)

		// Fill other game days with added player because on other days, he can go on the bench
		for day, pos := range pos_map {

				// Add player to gene for that day. When added here, doesn't count as a new player
				chromosome.Genes[day].Roster[pos] = player
		}
		cur_streamers[i] = player
	}
}


// Function to delete all occurences of a value in a chromosome
func DeleteAllOccurrences(chromosome *Chromosome, cur_streamers []Player, player_to_add Player, player_to_drop Player, week string, start_day int) {

	chromosome.DroppedPlayers[player_to_drop.Name] = DroppedPlayer{Player: player_to_drop, Countdown: 3}
	chromosome.Genes[start_day].DroppedPlayers = append(chromosome.Genes[start_day].DroppedPlayers, player_to_drop)

	for day := start_day; day <= ScheduleMap[week].GameSpan; day++ {

		// If the player doesn't have a game on the day, skip it and remove from bench
		if !Contains(ScheduleMap[week].Games[player_to_drop.Team], day) {
			chromosome.Genes[day].Bench = Remove(chromosome.Genes[day].Bench, SliceIndexOf(chromosome.Genes[day].Bench, player_to_drop))
			continue
		}

		key := MapContainsValue(chromosome.Genes[day].Roster, player_to_drop.Name)
		if key != "" {
			delete(chromosome.Genes[day].Roster, key)
		}
	}

	// Remove player from added players by replacing with new player
	index := SliceIndexOf(cur_streamers, player_to_drop)
	cur_streamers[index] = player_to_add
}

// Function to find the best position for a player
func FindBestPositions(player Player, chromosome *Chromosome, free_positions map[int][]string, pos_map map[int]string, start_day int, week string, add *bool) {

	// Loop through each day and find the best position for each day
	for day := start_day; day <= ScheduleMap[week].GameSpan; day++ {

		// If the player doesn't have a game on the day, skip it and add to bench
		if !Contains(ScheduleMap[week].Games[player.Team], day) {
			chromosome.Genes[day].Bench = append(chromosome.Genes[day].Bench, player)
			continue
		}

		has_match := false
		matches := GetMatches(player.ValidPositions, free_positions[day], &has_match)

		// If there are no matches, skip the day
		if !has_match {
			continue
		}

		// Go through matches in decreasing restriction order and assign the player to the first match that doesn't have a player in it
		for _, pos := range matches {
			
			// If the position doesn't have a player in it, add to pos_map and break
			if player, ok := chromosome.Genes[day].Roster[pos]; !ok || player.Name == "" {
				*add = true
				pos_map[day] = pos
				break
			}
		}

		// If we got here without adding a position, don't add the player
	}
}