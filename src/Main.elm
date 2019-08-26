module Main exposing (Model, Msg(..), initialModel, update, view)

import Browser
import Html exposing (..)
import Html.Attributes exposing (action, checked, class, id, method, name, placeholder, src, title, type_, value)
import Html.Events exposing (onClick, onInput)
import Http
import Process exposing (..)
import Task exposing (..)


type alias ResponseFromServer =
    { result : String
    }


type alias Model =
    { content : String
    }


makeFakeResponse : String -> ResponseFromServer
makeFakeResponse input =
    let
        result =
            input
            |> String.split " " 
            |> List.reverse 
            |> String.join " "
    in
    { result = result }


sendContentToServer : Model -> Cmd Msg
sendContentToServer model =
    let
        myTask =
            model.content
                |> makeFakeResponse
                |> Task.succeed
    in
    Process.sleep 1000
        |> Task.andThen (always myTask)
        |> Task.attempt HandleResponse


initialModel : Model
initialModel =
    { content = ""
    }


view : Model -> Html Msg
view model =
    div
        []
        [ h1 [] [ text "Echo Demo" ]
        , input
            [ type_ "text"
            , placeholder "write something..."
            , onInput ChangeContent
            , value model.content
            ]
            []
        , div [] [ button [ type_ "button", onClick SendContent ] [ text " OK " ] ]
        ]


type Msg
    = ChangeContent String
    | SendContent
    | HandleResponse (Result Http.Error ResponseFromServer)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        ChangeContent newContent ->
            ( { model | content = newContent }, Cmd.none )

        SendContent ->
            ( model, sendContentToServer model )

        HandleResponse (Ok response) ->
            ( { model | content = response.result }, Cmd.none )

        HandleResponse (Err _) ->
            ( model, Cmd.none )


main : Program () Model Msg
main =
    Browser.element
        { init = \_ -> ( initialModel, Cmd.none )
        , subscriptions = \_ -> Sub.none
        , update = update
        , view = view
        }
