package main

import "github.com/rs/zerolog/log"

func main() {

	err := Start()
	if err != nil {
		log.Info().Msgf("启动失败")
	}

}
