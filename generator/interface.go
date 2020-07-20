package generator

type PasswordGeneratorInterface interface {
	Generate() (string, bool)
}
