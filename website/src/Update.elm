module Update exposing (Msg(..), subscriptions, update)

import Model exposing (..)
import Process
import Task
import Translator exposing (..)


type Msg
    = ChangeTab Tab
    | UpdatedGameAppeared RawGame
    | LegalMovesAppeared (List Move)
    | MakeAction JsActions
    | ChangePlr1 String
    | ChangePlr2 String
    | NewGame
    | FlipBoard
    | UnselectStartI
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
        ChangeTab newTab ->
            ( { model | tab = newTab }, Cmd.none )

        UpdatedGameAppeared rg ->
            -- ! important
            -- TODO: verify that didn't change opponents in the mean time
            -- new game appeared, should make next move
            let
                lgd =
                    model.lgd

                -- so needs to figure out who that is that should make next move
                opponentToGo =
                    if rg.plrTurn == 0 then
                        lgd.plr1blue

                    else
                        lgd.plr2red

                updatedLgd =
                    { lgd | rg = rg, legalMoves = [], selectedStartI = Nothing }
            in
            -- reset the model
            ( { model | lgd = updatedLgd }
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
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | legalMoves = moves }
            in
            ( { model | lgd = updatedLgd }, Cmd.none )

        ChangePlr1 s ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | futurePlr1blue = getPlrSelected s }
            in
            ( { model | lgd = updatedLgd }, Cmd.none )

        ChangePlr2 s ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | futurePlr2red = getPlrSelected s }
            in
            ( { model | lgd = updatedLgd }, Cmd.none )

        NewGame ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | plr1blue = lgd.futurePlr1blue, plr2red = lgd.futurePlr2red }
            in
            -- flip future players to current
            -- start new game by calling update
            { model | lgd = updatedLgd }
                |> update (UpdatedGameAppeared startingRawGame)

        FlipBoard ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | boardFlipped = not lgd.boardFlipped }
            in
            ( { model | lgd = updatedLgd }, Cmd.none )

        StartSlotSelected i ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd
                        | selectedStartI =
                            if lgd.selectedStartI == Just i then
                                -- if wants to un-select
                                Nothing

                            else
                                Just i
                    }
            in
            ( { model | lgd = updatedLgd }
            , Cmd.none
            )

        UnselectStartI ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | selectedStartI = Nothing }
            in
            ( { model | lgd = updatedLgd }, Cmd.none )

        EndSlotSelected i ->
            let
                lgd =
                    model.lgd

                updatedLgd =
                    { lgd | selectedStartI = Nothing }
            in
            -- should only be possible to get here if model.selectedStartI is a value
            case lgd.selectedStartI of
                Nothing ->
                    ( model, Cmd.none )

                Just si ->
                    ( { model | lgd = updatedLgd }
                    , translator (MakeMove model.lgd.rg { startI = si, endI = i })
                    )


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ rawGameReceiver UpdatedGameAppeared
        , legalMovesReceiver LegalMovesAppeared
        ]
