package ini

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "strings"
)

func ParseFromFile(filepath string, data map[string]string) error {
    file, err := os.Open(filepath);
    if err != nil {
        return err
    }
    defer file.Close()

    const trimSet = " \"'."
    var prefix string
    lineNum := 0
    scanner := bufio.NewScanner(file)
    var key, val string
    for scanner.Scan() {
        line := scanner.Text()
        if len(line) == 0 {
            continue
        }
        if commentIdx := strings.Index(line, ";"); commentIdx != -1 {
            line = line[0:commentIdx]
            if commentIdx == 0 {
                continue
            }
        }
        if line[0] == '[' {
            if rabIdx := strings.Index(line, "]"); rabIdx != -1 {
                prefix = line[1:rabIdx]
                continue
            }
            return errors.New(fmt.Sprintf("Syntax error in %s at line %d: incomplete section declaration.", filepath, lineNum))
        }
        if lsbIdx := strings.Index(line, "["); lsbIdx != -1 {
            return errors.New(fmt.Sprintf("Syntax error in %s at line %d: unexpected [; only use [ and ] for section declaration.", filepath, lineNum))
        }
        if rsbIdx := strings.Index(line, "]"); rsbIdx != -1 {
            return errors.New(fmt.Sprintf("Syntax error in %s at line %d: unexpected ]; only use [ and ] for section declaration.", filepath, lineNum))
        }
        if eqIdx := strings.Index(line, "="); eqIdx != -1 {
            key = line[:eqIdx]
            val = line[eqIdx+1:]
        } else {
            return errors.New(fmt.Sprintf("Syntax error in %s at line %d: no '=' found.", filepath, lineNum))
        }
        if pIdxK := strings.Index(key, "."); pIdxK != -1 {
            return errors.New(fmt.Sprintf("Syntax error in %s at line %d: illegal usage of '.' in key name.", filepath, lineNum))
        }

        data[fmt.Sprintf("%s.%s", prefix, strings.Trim(key, trimSet))] = strings.Trim(val, trimSet)
        lineNum++
    }

    return nil
}

func WriteToFile(filepath string, data map[string]string) error {
    file, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    section := ""
    for k, v := range data {
        if !strings.Contains(k, section) || section == "" {
            if !strings.Contains(k, ".") {
                return errors.New(fmt.Sprintf("In function ini.WriteToFile; param `data map[string]string`: key %s is not associated a section.", k))
            }
        }
        if secDelIdx := strings.Index(k, "."); secDelIdx != -1 {
            kSection := k[:secDelIdx]
            k = k[secDelIdx+1:]
            if kSection != section {
                section = kSection
                if _, err := file.WriteString(fmt.Sprintf("[%s]\n", kSection)); err != nil {
                    return err
                }
            }
        } else {
            return errors.New(fmt.Sprintf("In function ini.WriteToFile; param `data map[string]string`: key %s is not associated a section.", k))
        }

        if _, err := file.WriteString(fmt.Sprintf("%s = %s\n", k, v)); err != nil {
            return err
        }
    }

    return nil
}
