module Update exposing (..)

import Model exposing (..)
import Translator exposing (..)


type Msg
    = UpdatedGameAppeared RawGame
    | LegalMovesAppeared (List Move)
    | MakeAction JsActions
    | ChangeDifficulty String
    | NewGame
    | SlotSelected Int


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdatedGameAppeared rg ->
            ( { model | rg = rg }, translator (GetLegalMoves rg) )

        LegalMovesAppeared moves ->
            ( { model | legalMoves = moves, selectedStartI = Nothing }, Cmd.none )

        MakeAction action ->
            ( model, translator action )

        ChangeDifficulty stringedDiff ->
            ( { model
                | difficulty =
                    case stringedDiff of
                        "Easy" ->
                            Easy

                        "Medium" ->
                            Medium

                        "Hard" ->
                            Hard

                        "ExtraHard" ->
                            ExtraHard

                        "Impossible" ->
                            Impossible

                        "Simple" ->
                            Simple

                        _ ->
                            Simple
              }
            , Cmd.none
            )

        NewGame ->
            ( { model
                | rg = startingRawGame
                , legalMoves = startingLegalMoves
              }
            , Cmd.none
            )

        SlotSelected i ->
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


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ rawGameReceiver UpdatedGameAppeared
        , legalMovesReceiver LegalMovesAppeared
        ]
