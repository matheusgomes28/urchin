function HandleShortcode(arguments)
    if #arguments ~= 2 then
        return ""
    end

    return "\n| Tables   |      Are      |  Cool |\n|----------|:-------------:|------:|\n| col 1 is |  left-aligned | $1600 |\n| col 2 is |    centered   |   $12 |\n| col 3 is | right-aligned |    $1 |\n"
end
