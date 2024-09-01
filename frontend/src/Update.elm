module Update exposing (..)

import Model exposing (..)
import Process
import Task
import Translator exposing (..)


type Msg
    = UpdatedGameAppeared RawGame
    | LegalMovesAppeared (List Move)
    | MakeAction JsActions
    | ChangePlr1 String
    | ChangePlr2 String
    | NewGame
    | FlipBoard
    | StartSlotSelected Int
    | EndSlotSelected Int


getPlrSelected : String -> Opponent
getPlrSelected s =
    case s of
        "Human" ->
            Human

        "Easy" ->
            Ai Easy

        "Medium" ->
            Ai Medium

        "Hard" ->
            Ai Hard

        "ExtraHard" ->
            Ai ExtraHard

        "Impossible" ->
            Ai Impossible

        "Simple" ->
            Ai Simple

        _ ->
            Human


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdatedGameAppeared rg ->
            -- new game appeared, should make next move
            let
                -- so needs to figure out who that is that should make next move
                opponentToGo =
                    if rg.plrTurn == 0 then
                        model.plr1blue

                    else
                        model.plr2red
            in
            -- reset the model
            ( { model | rg = rg, legalMoves = [], selectedStartI = Nothing }
            , if rg.state == "Playing" then
                case opponentToGo of
                    Human ->
                        -- if human should go, then should get legal moves so it
                        -- will be displayed and clickable
                        translator (GetLegalMoves rg)

                    Ai diff ->
                        -- not a great solution to make the board display and then get move
                        Process.sleep 300
                            |> Task.perform
                                (always (MakeAction <| GetAiMove rg diff))

              else
                Cmd.none
            )

        MakeAction action ->
            ( model, translator action )

        LegalMovesAppeared moves ->
            ( { model | legalMoves = moves }, Cmd.none )

        ChangePlr1 s ->
            ( { model | futurePlr1blue = getPlrSelected s }, Cmd.none )

        ChangePlr2 s ->
            ( { model | futurePlr2red = getPlrSelected s }, Cmd.none )

        NewGame ->
            -- flip future players to current
            -- start new game by calling update
            { model | plr1blue = model.futurePlr1blue, plr2red = model.futurePlr2red }
                |> update (UpdatedGameAppeared startingRawGame)

        FlipBoard ->
            ( { model | boardFlipped = not model.boardFlipped |> Debug.log "boardFlipped" }, Cmd.none )

        StartSlotSelected i ->
            ( { model
                | selectedStartI =
                    if model.selectedStartI == Just i then
                        -- if wants to un-select
                        Nothing

                    else
                        Just i
              }
            , Cmd.none
            )

        EndSlotSelected i ->
            -- should only be possible to get here if model.selectedStartI is a value
            case model.selectedStartI of
                Nothing ->
                    ( model, Cmd.none )

                Just si ->
                    ( { model | selectedStartI = Nothing }
                    , translator (MakeMove model.rg { startI = si, endI = i })
                    )


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ rawGameReceiver UpdatedGameAppeared
        , legalMovesReceiver LegalMovesAppeared
        ]
