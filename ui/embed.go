package ui

import "embed"

// DistFS expõe os arquivos da pasta 'dist' que serão embutidos no binário executável final do Go.
//
// A diretiva '//go:embed dist/*' é um comando especial de pré-compilação para o Go.
// Ela instrui o compilador a ler todo o conteúdo do diretório 'dist' e guardá-lo
// diretamente na memória do executável final como uma variável do tipo 'embed.FS' (Filesystem).
//
// Analogia com Flutter/Dart:
// No Flutter, para usar imagens ou arquivos HTML locais, declaramos eles no arquivo 'pubspec.yaml':
//   flutter:
//     assets:
//       - ui/dist/
// E então usamos 'rootBundle.loadString()' para lê-los.
// Em Go, fazemos isso declarando a diretiva '//go:embed' logo acima da variável que servirá como
// o nosso sistema de arquivos virtual embutido.
//
//go:embed dist/*
var DistFS embed.FS
