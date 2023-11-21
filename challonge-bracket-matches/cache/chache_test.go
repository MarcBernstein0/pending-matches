package cache

import (
	"testing"
	"time"

	"github.com/MarcBernstein0/pending-matches/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCache(t *testing.T) {
	mockCache := NewCache(5*time.Minute, 5*time.Hour)
	// Given
	tt := []struct {
		name     string
		mockData []models.TournamentParticipants
	}{
		{
			name: "Single Tournament",
			mockData: []models.TournamentParticipants{
				{
					GameName:     "Guilty Gear -Strive-",
					TournamentID: 10879090,
					Participant: map[int]string{
						166014671: "test",
						166014672: "test2",
						166014673: "test3",
						166014674: "test4",
					},
				},
			},
		},
		{
			name: "Multiple tournaments",
			mockData: []models.TournamentParticipants{
				{
					GameName:     "Guilty Gear -Strive-",
					TournamentID: 10879090,
					Participant: map[int]string{
						166014671: "test",
						166014672: "test2",
						166014673: "test3",
						166014674: "test4",
					},
				},
				{
					GameName:     "DNF Duel",
					TournamentID: 10879091,
					Participant: map[int]string{
						166014671: "test",
						166014672: "test2",
						166014673: "test3",
						166014674: "test4",
					},
				},
			},
		},
	}
	for _, tC := range tt {
		t.Run(tC.name, func(t *testing.T) {
			// When
			mockCache.UpdateCache(tC.mockData, "2006-01-02")
			// Then
			assert.ElementsMatch(t, tC.mockData, mockCache.GetData("2006-01-02"))
		})

	}
}

func TestShouldUpdate(t *testing.T) {
	t.Run("Test that method returns true when timer has exceeded limit", func(t *testing.T) {
		// Given
		mockCache := NewCache(2*time.Microsecond, 2*time.Microsecond)
		mockCache.UpdateCache([]models.TournamentParticipants{
			{
				GameName:     "Guilty Gear -Strive-",
				TournamentID: 10879090,
				Participant: map[int]string{
					166014671: "test",
					166014672: "test2",
					166014673: "test3",
					166014674: "test4",
				},
			},
		}, "2006-01-02")
		// When
		time.Sleep(5 * time.Millisecond)
		// Then
		assert.Equal(t, true, mockCache.ShouldUpdate("2006-01-02"))
	})

	t.Run("Test that method returns false when timer has note exceeded limit", func(t *testing.T) {
		// Given
		mockCache := NewCache(5*time.Millisecond, 5*time.Millisecond)
		mockCache.UpdateCache([]models.TournamentParticipants{
			{
				GameName:     "Guilty Gear -Strive-",
				TournamentID: 10879090,
				Participant: map[int]string{
					166014671: "test",
					166014672: "test2",
					166014673: "test3",
					166014674: "test4",
				},
			},
		}, "2006-01-02")
		// When
		time.Sleep(2 * time.Microsecond)
		// Then
		assert.Equal(t, false, mockCache.ShouldUpdate("2006-01-02"))
	})
}
