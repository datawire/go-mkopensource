package main

type exitCode int

const (
	NoError                   exitCode = 0
	DependencyGenerationError exitCode = 1
	MarshallJsonError         exitCode = 2
	WriteError                exitCode = 3
	InvalidArgumentsError     exitCode = 4
)
