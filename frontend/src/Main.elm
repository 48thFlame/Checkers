module Main exposing (main)

import Browser
import Html
import Html.Attributes exposing (alt, class, src, style)


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
    | BluePiece
    | BlueKing
    | RedPiece
    | RedKing


newBoard : List Slot
newBoard =
    [ [ NaS, BluePiece, NaS, BluePiece, NaS, BluePiece, NaS, BluePiece ]
    , [ BluePiece, NaS, BluePiece, NaS, BluePiece, NaS, BluePiece, NaS ]
    , [ NaS, BluePiece, NaS, BluePiece, NaS, BluePiece, NaS, BluePiece ]
    , [ Empty, NaS, Empty, NaS, BlueKing, NaS, Empty, NaS ]
    , [ NaS, Empty, NaS, RedKing, NaS, Empty, NaS, Empty ]
    , [ RedPiece, NaS, RedPiece, NaS, RedPiece, NaS, RedPiece, NaS ]
    , [ NaS, RedPiece, NaS, RedPiece, NaS, RedPiece, NaS, RedPiece ]
    , [ RedPiece, NaS, RedPiece, NaS, RedPiece, NaS, RedPiece, NaS ]
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


slotToHtml : Int -> Slot -> Html.Html msg
slotToHtml i slot =
    let
        slotISpan =
            Html.span [ class "slot-i" ]
                [ Html.text <| String.fromInt i ]

        slotImg pieceName =
            Html.img
                [ src ("assets/" ++ pieceName ++ ".webp")
                , alt pieceName
                , class "slot-img"
                ]
                []
    in
    case slot of
        NaS ->
            Html.div [ class "slot NaS-slot" ] []

        Empty ->
            Html.div [ class "slot empty-slot" ]
                [ slotISpan ]

        BluePiece ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "bluePiece"
                ]

        BlueKing ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "blueKing"
                ]

        RedPiece ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "redPiece"
                ]

        RedKing ->
            Html.div [ class "slot piece-slot" ]
                [ slotISpan
                , slotImg "redKing"
                ]


boardToHtml : List Slot -> Html.Html Msg
boardToHtml board =
    let
        stringedBoardSize =
            String.fromInt 8
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        , style "grid-template-columns" ("repeat(" ++ stringedBoardSize ++ ", 1fr)")
        ]
        (List.indexedMap slotToHtml board)


view : Model -> Html.Html Msg
view model =
    Html.div [ class "main-area" ]
        [ boardToHtml model.board ]
