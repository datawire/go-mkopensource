package main

type exitCode int

const (
	DependencyGenerationError exitCode = 1
	MarshallJsonError         exitCode = 2
	WriteError                exitCode = 3
)
