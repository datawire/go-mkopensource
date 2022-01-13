package main

type exitCode int

const (
	ShowProgramHelp exitCode = iota
	DependencyGenerationError
	InvalidArgumentsError
	MarshallJsonError
)
