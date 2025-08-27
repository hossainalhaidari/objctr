export OBJCTR_HOST="<OBJCTR_HOST>"
export OBJCTR_AUTH="Authorization: <OBJCTR_AUTH>"

objctr () {
    case $1 in
        "ls")
            curl $OBJCTR_HOST/$2 -H "$OBJCTR_AUTH"
            ;;
        "mk")
            if [ ! -f "$3" ]; then
                curl $OBJCTR_HOST/$2 -X POST -H "$OBJCTR_AUTH"
            else
                curl $OBJCTR_HOST/$2 -X POST -H "$OBJCTR_AUTH" -F "file=@$3"
            fi
            ;;
        "cp")
            curl $OBJCTR_HOST/$2?to=$3 -X PUT -H "$OBJCTR_AUTH"
            ;;
        "mv")
            curl $OBJCTR_HOST/$2?to=$3 -X PATCH -H "$OBJCTR_AUTH"
            ;;
        "rm")
            curl $OBJCTR_HOST/$2 -X DELETE -H "$OBJCTR_AUTH"
            ;;
        *)
            echo "Unknown command: $1"
            ;;
    esac
}
