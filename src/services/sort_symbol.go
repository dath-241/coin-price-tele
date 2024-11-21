package services

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"path/filepath"
)

type SymbolGroup struct {
	usdt  []string
	usdc  []string
	btc   []string
	other []string
}

func CallSortSymbols() error {
    fmt.Println("Sorting symbols...")

	// Lấy đường dẫn của thư mục hiện tại
    currentDir, err := os.Getwd()
    if err != nil {
        return fmt.Errorf("error getting current directory: %v", err)
    }

    // Define base directory for files
    files := []string{
        filepath.Join(currentDir, "services", "spot_symbols.txt"),
        filepath.Join(currentDir, "services", "futures_symbols.txt"),
    }

    // Check if files exist with detailed logging
    for _, file := range files {
        fmt.Printf("Checking file: %s\n", file)
        if _, err := os.Stat(file); os.IsNotExist(err) {
            return fmt.Errorf("file %s does not exist", file)
        }
        fmt.Printf("File %s exists\n", file)
    }

    // Sort spot symbols with updated paths
	fmt.Println("Sorting spot symbols...")
    if err := SortSymbols(
        filepath.Join(currentDir,  "services", "spot_symbols.txt"),
        filepath.Join(currentDir,  "services", "spot_symbols_sorted.txt"),
    ); err != nil {
        return fmt.Errorf("error sorting spot symbols: %v", err)
    }

    // Sort futures symbols with updated paths
    fmt.Println("Sorting futures symbols...")
    if err := SortSymbols(
        filepath.Join(currentDir, "services", "futures_symbols.txt"),
        filepath.Join(currentDir, "services", "futures_symbols_sorted.txt"),
    ); err != nil {
        return fmt.Errorf("error sorting futures symbols: %v", err)
    }

    fmt.Println("Symbols sorted successfully!")
    return nil
}

func SortSymbols(inputFile, outputFile string) error {
    // Đọc file
    file, err := os.Open(inputFile)
    if err != nil {
        return err
    }
    defer file.Close()

    // Đọc các dòng từ file
    scanner := bufio.NewScanner(file)
    var header []string
    var symbols []string
    lineCount := 0

    for scanner.Scan() {
        if lineCount < 2 {
            header = append(header, scanner.Text()+"\n")
        } else {
            line := strings.TrimSpace(scanner.Text())
            if line != "" {
                symbols = append(symbols, line)
            }
        }
        lineCount++
    }

    // Tạo map để nhóm các cặp theo tiền tố
    prefixGroups := make(map[string]*SymbolGroup)

    for _, symbol := range symbols {
        var prefix string

        // Tìm tiền tố
        switch {
        case strings.HasSuffix(symbol, "USDT"):
            prefix = symbol[:len(symbol)-4]
        case strings.HasSuffix(symbol, "USDC"):
            prefix = symbol[:len(symbol)-4]
        case strings.HasSuffix(symbol, "BTC"):
            prefix = symbol[:len(symbol)-3]
        default:
            prefix = symbol
        }

        // Khởi tạo group nếu chưa tồn tại
        if _, exists := prefixGroups[prefix]; !exists {
            prefixGroups[prefix] = &SymbolGroup{}
        }

        // Phân loại symbol
        switch {
        case strings.HasSuffix(symbol, "USDT"):
            prefixGroups[prefix].usdt = append(prefixGroups[prefix].usdt, symbol)
        case strings.HasSuffix(symbol, "USDC"):
            prefixGroups[prefix].usdc = append(prefixGroups[prefix].usdc, symbol)
        case strings.HasSuffix(symbol, "BTC"):
            prefixGroups[prefix].btc = append(prefixGroups[prefix].btc, symbol)
        default:
            prefixGroups[prefix].other = append(prefixGroups[prefix].other, symbol)
        }
    }

    // Sắp xếp các symbol trong mỗi group theo thứ tự alphabet
    for _, group := range prefixGroups {
        sort.Strings(group.usdt)
        sort.Strings(group.usdc)
        sort.Strings(group.btc)
        sort.Strings(group.other)
    }

    // Ghi ra file mới
    outFile, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer outFile.Close()

    writer := bufio.NewWriter(outFile)

    // Ghi header
    for _, h := range header {
        writer.WriteString(h)
    }

    // Sắp xếp prefixes
    prefixes := make([]string, 0, len(prefixGroups))
    for prefix := range prefixGroups {
        prefixes = append(prefixes, prefix)
    }
    sort.Strings(prefixes)

    // Định nghĩa thứ tự ưu tiên của postfix
    postfixOrder := []struct {
        symbols func(*SymbolGroup) []string
    }{
        {func(g *SymbolGroup) []string { return g.usdt }},
        {func(g *SymbolGroup) []string { return g.usdc }},
        {func(g *SymbolGroup) []string { return g.btc }},
        {func(g *SymbolGroup) []string { return g.other }},
    }

    // Ghi theo thứ tự prefix và postfix
    for _, prefix := range prefixes {
        group := prefixGroups[prefix]
        for _, getter := range postfixOrder {
            symbols := getter.symbols(group)
            if len(symbols) > 0 {
                writer.WriteString(strings.Join(symbols, "\n") + "\n")
            }
        }
    }

    return writer.Flush()
}