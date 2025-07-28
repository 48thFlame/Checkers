module View.Board exposing (viewBoard)

import Html
import Html.Attributes exposing (alt, class, draggable, src, style)
import Html.Events exposing (onMouseDown, onMouseUp)
import Model exposing (Model, Move)
import Set exposing (Set)
import Update exposing (Msg(..))


type PieceSlot
    = NaS
    | Empty
    | BluePiece
    | BlueKing
    | RedPiece
    | RedKing


type SlotFlags
    = MoveStater
    | Selected
    | MoveEnder
    | Regular


type alias BoardSlot =
    ( PieceSlot, SlotFlags )


{-| Convert a raw game board to a list of `PieceSlot`
-}
rgBoardToPieces : List Int -> List PieceSlot
rgBoardToPieces board =
    let
        intToPiece : Int -> PieceSlot
        intToPiece n =
            case n of
                0 ->
                    NaS

                1 ->
                    Empty

                2 ->
                    BluePiece

                3 ->
                    BlueKing

                4 ->
                    RedPiece

                5 ->
                    RedKing

                _ ->
                    NaS
    in
    List.map intToPiece board


{-| Convert a raw game board to a list of `SlotFlags`
-}
rgBoardToSlotFlags : Maybe Int -> List Move -> List Int -> List SlotFlags
rgBoardToSlotFlags selectedStartI legalMoves board =
    let
        {- Get a Set of all legal moves starting Is
           so it can become a `MoveStarter`
        -}
        startIs : Set Int
        startIs =
            List.map (\move -> move.startI) legalMoves
                |> Set.fromList

        {- Get a list of all legal moves ending Is
           if user has a slot selcted so can show all the `MoveEnder`s
        -}
        endIs : List Int
        endIs =
            case selectedStartI of
                Nothing ->
                    []

                Just si ->
                    List.filter (\move -> move.startI == si) legalMoves
                        |> List.map (\move -> move.endI)
    in
    List.indexedMap
        (\i _ ->
            if Set.member i startIs then
                if selectedStartI == Just i then
                    Selected

                else
                    MoveStater

            else if List.member i endIs then
                MoveEnder

            else
                Regular
        )
        board


pieceSlotToHtml : Int -> PieceSlot -> Html.Html Msg
pieceSlotToHtml i pieceSlot =
    let
        slotISpan : Html.Html Msg
        slotISpan =
            Html.span [ class "slot-i-span-indicator" ]
                [ Html.text <| String.fromInt i ]

        slotImg : String -> Html.Html Msg
        slotImg pieceName =
            Html.img
                [ src ("assets/" ++ pieceName ++ ".webp")
                , alt pieceName
                , class "piece-img-element"
                , draggable "false"
                ]
                []

        slot : String -> Html.Html Msg
        slot name =
            Html.div [ class "inside-slot-house-of-piece" ]
                [ slotISpan
                , slotImg name
                ]
    in
    case pieceSlot of
        NaS ->
            Html.div [ class "NaS-type-slot" ] []

        Empty ->
            -- Html.div [ class "empty-type-slot" ] [ slotISpan ]
            Html.div [] [ slotISpan ]

        BluePiece ->
            slot "bluePiece"

        BlueKing ->
            slot "blueKing"

        RedPiece ->
            slot "redPiece"

        RedKing ->
            slot "redKing"


viewBoardSlot : Bool -> Int -> BoardSlot -> Html.Html Msg
viewBoardSlot isBoardFlipped i ( piece, flag ) =
    let
        rowStr : String
        rowStr =
            (i // 8)
                + 1
                |> (if isBoardFlipped then
                        \rn -> 9 - rn

                    else
                        identity
                   )
                |> String.fromInt

        colStr : String
        colStr =
            modBy 8 i
                + 1
                |> (if isBoardFlipped then
                        \cl -> 9 - cl

                    else
                        identity
                   )
                |> String.fromInt

        baseAttrs : List (Html.Attribute Msg)
        baseAttrs =
            [ style "grid-row-start" rowStr
            , style "grid-column-start" colStr
            , class "outer-slot"
            ]

        deSelectStartIEvent : Html.Attribute Msg
        deSelectStartIEvent =
            onMouseUp UnselectStartI

        pieceSlot : Html.Html Msg
        pieceSlot =
            pieceSlotToHtml i piece
    in
    case flag of
        Regular ->
            Html.div
                (deSelectStartIEvent :: baseAttrs)
                [ pieceSlot ]

        MoveStater ->
            Html.div
                ([ class "starting-legal-move-slot"
                 , onMouseDown (StartSlotSelected i)
                 ]
                    ++ baseAttrs
                )
                [ pieceSlot ]

        Selected ->
            Html.div
                ([ class "starting-legal-move-slot selected-slot-to-make-move", deSelectStartIEvent
                ]
                    ++ baseAttrs
                )
                [ pieceSlot ]

        MoveEnder ->
            Html.div
                ([ class "ending-legal-move-slot", onMouseUp (EndSlotSelected i) ] ++ baseAttrs)
                [ Html.div [ class "ending-legal-move-circle-indicator" ] []
                , pieceSlot
                ]


viewBoard : Model -> Html.Html Msg
viewBoard model =
    let
        pieces : List PieceSlot
        pieces =
            rgBoardToPieces model.rg.board

        slotFlags : List SlotFlags
        slotFlags =
            rgBoardToSlotFlags model.selectedStartI model.legalMoves model.rg.board

        boardSlots : List BoardSlot
        boardSlots =
            List.map2 Tuple.pair pieces slotFlags

        viewSlots : List (Html.Html Msg)
        viewSlots =
            List.indexedMap (viewBoardSlot model.boardFlipped) boardSlots

        gameStateDiv : Html.Html Msg
        gameStateDiv =
            Html.div [ class "game-state", style "grid-area" "1/1/9/9" ]
                [ Html.p [ class "game-state-text" ]
                    [ Html.text
                        (case model.rg.state of
                            "Draw" ->
                                "Game ended in a draw.."

                            "Blue Won" ->
                                "Blue won the game!!"

                            "Red Won" ->
                                "Red won the game!!"

                            _ ->
                                ""
                        )
                    , Html.br [] []
                    , Html.br [] []
                    , Html.text "Press \"Play\" button to play another game."
                    ]
                ]
    in
    Html.div
        [ class "checker-board"
        , style "grid-template-rows" "repeat(8, 1fr)"
        , style "grid-template-columns" "repeat(8, 1fr)"
        ,  style "cursor" (if model.selectedStartI /= Nothing then
             "grabbing"
         else "default"
           )
            
        ]
        (viewSlots
            |> (if model.rg.state /= "Playing" then
                    List.append [ gameStateDiv ]

                else
                    identity
               )
        )
