package hardware

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/hardware/sequent"
	"github.com/d2r2/go-i2c"
	"go.uber.org/zap"
)

type relay8Board struct {
	dev   *i2c.I2C
	stack byte
	state sequent.Relay8States
}

/*func main() {
	var relay8s []relay8Board
	for i := byte(0); i < 8; i++ {
		dev, err := sequent.InitBoard(i)
		if err != nil {
			continue
		}

		relay8s = append(relay8s, relay8Board{
			dev:   dev,
			stack: i,
			state: sequent.Relay8States{},
		})
	}

	if len(relay8s) == 0 {
		log.Fatal("no relay8 boards found")
	}

	defer func() {
		for i := range relay8s {
			board := relay8s[i]
			sequent.DeinitBoard(board.dev)
			board.dev = nil
		}
	}()

	for i := 0; i < 256; i++ {
		for j := range relay8s {
			board := relay8s[j]
			initialState := board.state
			for k := 0; k < 8; k++ {
				board.state = board.state.Set(k, byte(i&(1<<k)) != 0)
			}

			if !initialState.Equal(board.state) {
				err := sequent.UpdateBoard(board.dev, board.state, 10)
				if err != nil {
					log.Println(fmt.Errorf("error updating board %d: %w", board.stack, err))
				}
			}
		}
		time.Sleep(time.Second)
	}
}*/

type SequentRelay8Hardware struct {
	mu     *sync.Mutex
	boards []relay8Board
}

func NewSR8Hardware() (*SequentRelay8Hardware, error) {
	var relay8s []relay8Board
	for i := byte(0); i < 8; i++ {
		dev, err := sequent.InitBoard(i)
		if err != nil {
			continue
		}

		relay8s = append(relay8s, relay8Board{
			dev:   dev,
			stack: i,
			state: sequent.Relay8States{},
		})
	}

	if len(relay8s) == 0 {
		return nil, fmt.Errorf("no relay8 boards found")
	}

	return &SequentRelay8Hardware{
		mu:     &sync.Mutex{},
		boards: relay8s,
	}, nil
}

func (s *SequentRelay8Hardware) Name() string {
	return "sequent-relay8"
}

func (s *SequentRelay8Hardware) Close() error {
	for i := range s.boards {
		board := s.boards[i]
		sequent.DeinitBoard(board.dev)
		board.dev = nil
	}

	return nil
}

func (s *SequentRelay8Hardware) NumPumps() int {
	return len(s.boards) * 8
}

func (s *SequentRelay8Hardware) Pump(idx int, state PumpState) error {
	if idx < 0 || idx >= s.NumPumps() {
		return fmt.Errorf("invalid pump index %d", idx)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	boardIdx := idx / 8
	pumpIdx := idx % 8

	board := s.boards[boardIdx]
	board.state = board.state.Set(pumpIdx, state != Off)

	return nil
}

func (s *SequentRelay8Hardware) Update(logger *zap.Logger) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.boards {
		board := s.boards[i]
		err := sequent.UpdateBoard(board.dev, board.state, 10)
		if err != nil {
			logger.Info("error updating board", zap.Uint8("board_idx", board.stack), zap.Error(err))
		}
	}
}
