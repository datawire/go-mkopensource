#!/bin/sh
set -e

JS_LICENSES="/temp/js_licenses.txt"
JS_DEPENDENCIES="/temp/js_dependencies.txt"

PY_LICENSES="/temp/py_licenses.txt"
PY_DEPENDENCIES="/temp/py_dependencies.txt"

GO_DEPENDENCIES="/temp/go_dependencies.txt"
GO_LICENSES="/temp/go_licenses.txt"

generate_opensource() {
  TMP_LICENSES=/tmp/licenses.txt

  INPUT="$1"
  LANGUAGE="$2"
  OUTPUT="$3"

  if [[ -f "${INPUT}" ]]; then
    {
      echo -e "Name|Version|License(s)\n----|-------|----------"
      cat "${INPUT}"
    } >"${TMP_LICENSES}"

    {
      echo -e "The ${APPLICATION} ${LANGUAGE} code makes use of the following Free and Open Source\nlibraries:\n"

      gawk 'BEGIN{OFS=FS="|"}
             NR==FNR {for (i=1; i<=NF; i++) max[i]=(length($i)>max[i]?length($i):max[i]); next}
                     {for (i=1; i<=NF; i++) printf "%s%-*s%s", i==1 ? "    " : "", i < NF? max[i]+2 : 1, $i, i==NF ? ORS : " "}
           ' "${TMP_LICENSES}" "${TMP_LICENSES}"
    } >"${OUTPUT}"
  else
    echo "File ${INPUT} does not exist. Skipping license generation from it" >&2
  fi
}

validate_required_variable() {
  VARIABLE="$1"

  if [[ -z "${!VARIABLE}" ]]; then
    echo "Variable ${VARIABLE} is required" >&2
    exit 1
  fi
}