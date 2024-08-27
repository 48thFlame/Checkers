module Update exposing (..)

import Model exposing (..)
import Translator exposing (..)


type Msg
    = UpdatedGameAppeared RawGame
    | LegalMovesAppeared (List Move)
    | MakeAction JsActions
    | ChangeDifficulty AiDifficulty


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdatedGameAppeared rg ->
            ( { model | rg = rg }, translator (GetLegalMoves rg) )

        LegalMovesAppeared moves ->
            ( { model | legalMoves = moves |> Debug.log "from Elm" }, Cmd.none )

        MakeAction action ->
            ( model, translator action )

        ChangeDifficulty diff ->
            ( { model | difficulty = diff }, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.batch
        [ rawGameReceiver UpdatedGameAppeared
        , legalMovesReceiver LegalMovesAppeared
        ]
