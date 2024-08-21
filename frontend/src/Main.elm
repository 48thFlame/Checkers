module Main exposing (main)

import Browser
import Html
import Html.Attributes exposing (class, style)


main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


type Slot
    = NaS
    | Empty



-- | BluePiece
-- | BlueKing
-- | RedPiece
-- | RedKing


newBoard : List Slot
newBoard =
    [ [ NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty ]
    , [ Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS ]
    , [ NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty ]
    , [ Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS ]
    , [ NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty ]
    , [ Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS ]
    , [ NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty ]
    , [ Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS ]
    ]
        |> List.concat


type alias Model =
    { board : List Slot
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( Model newBoard, Cmd.none )


type Msg
    = Msg1


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Msg1 ->
            ( model, Cmd.none )


subscriptions : Model -> Sub Msg
subscriptions _ =
    Sub.none


cellToHtml : Int -> Slot -> Html.Html msg
cellToHtml i cell =
    let
        stringedI =
            String.fromInt i
    in
    case cell of
        NaS ->
            Html.div [ class "cell NaS-cell" ] []

        Empty ->
            Html.div [ class "cell empty-cell" ] [ Html.text stringedI ]


boardToHtml : List Slot -> Html.Html Msg
boardToHtml board =
    let
        stringedSlotRow : Int -> String
        stringedSlotRow i =
            (i // 8) + 1 |> String.fromInt

        stringedSlotCol : Int -> String
        stringedSlotCol i =
            modBy 8 i + 1 |> String.fromInt

        slotInGridToHtml : Int -> Slot -> Html.Html Msg
        slotInGridToHtml i cell =
            Html.div
                [ style "grid-row-start" (stringedSlotRow i)
                , style "grid-column-start" (stringedSlotCol i)
                ]
                [ cellToHtml i cell ]

        stringedBoardSize =
            String.fromInt 8
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        , style "grid-template-columns" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        ]
        (List.indexedMap slotInGridToHtml board)


view : Model -> Html.Html Msg
view model =
    Html.div [ class "main-area" ]
        [ boardToHtml model.board ]
